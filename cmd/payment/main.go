package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vantu-fit/saga-pattern/cmd/payment/config"
	db "github.com/vantu-fit/saga-pattern/internal/payment/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/payment/event"
	"github.com/vantu-fit/saga-pattern/internal/payment/grpc"
	"github.com/vantu-fit/saga-pattern/internal/payment/http"
	"github.com/vantu-fit/saga-pattern/internal/payment/service"
	grpcclient "github.com/vantu-fit/saga-pattern/pkg/grpc_client"
	kafkaClient "github.com/vantu-fit/saga-pattern/pkg/kafka"
	migrate_db "github.com/vantu-fit/saga-pattern/pkg/migrate"
)

var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	// log for development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// start service
	log.Info().Msg("Start payment service")

	// load config
	cfgFile, err := config.LoadConfig("./cmd/payment/config/config")
	if err != nil {
		log.Fatal().Msgf("Load config: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatal().Msgf("Parse config: %v", err)
	}

	// create context for gracefull shutdown
	doneCh := make(chan struct{}) // for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), interuptSignals...)
	defer stop()

	// run migrate DB
	poolConfig, err := pgxpool.ParseConfig(cfg.Postgres.DnsURL)
	if err != nil {
		log.Fatal().Msgf("Parse pgx pool config: %v", err)
	}
	poolConfig.MaxConns = 500// Set maximum connections in the pool
	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	migrate_db.RunDBMigration(cfg.Postgres.Migration, cfg.Postgres.DnsURL)

	store := db.NewStore(conn)

	// create grpc client
	grpcClient := grpcclient.NewClient()

	// create order service
	paymentService := service.NewService(store)

	// create grpc server
	grpcServer := grpc.NewServer(cfg, store, paymentService, grpcClient)

	// run grpc server
	go func() {
		log.Info().Msgf("GRPC server is running on port %s", cfg.GRPC.Port)
		if err := grpcServer.Run(); err != nil {
			log.Fatal().Msgf("Run grpc server: %v", err)
		}
	}()

	// create http gateway server
	httpServer, err := http.NewHTTPGatewayServer(cfg, store, paymentService, grpcClient)
	if err != nil {
		log.Fatal().Msgf("Create HTTPGateway Server: %s", err)
	}
	// run http gateway server
	go func() {
		log.Info().Msgf("HTTPGateway Server is run on port %s", cfg.HTTP.Port)
		if err := httpServer.Run(); err != nil {
			log.Fatal().Msgf("Run HTTPGateway Server: %s", err)
		}
	}()

	// create kafka
	producer := kafkaClient.NewProducer(cfg.Kafka.Brokers)
	consumer := kafkaClient.NewConsumerGroup(cfg.Kafka.Brokers)

	// create event handler
	orderEvenHandler := event.NewEventHandler(cfg, consumer, producer, paymentService)

	// run event handler
	orderEvenHandler.Run(ctx)

	// graceful shutdown
	<-ctx.Done()
	go func() {
		time.Sleep(0 * time.Second)
		log.Fatal().Msg("Graceful shutdown timeout")

		// HTTPGatewayServer.Shutdown(context.Background())

		grpcServer.GracefulStop()
		doneCh <- struct{}{}
	}()

	log.Info().Msg("product service shutdown")
	<-doneCh
}
