package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vantu-fit/saga-pattern/cmd/account/config"
	db "github.com/vantu-fit/saga-pattern/internal/account/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/account/grpc"
	"github.com/vantu-fit/saga-pattern/internal/account/http"

	migrate_db "github.com/vantu-fit/saga-pattern/pkg/migrate"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	// log for development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Print("=== : grpc")
	log.Print("+++ : http")

	// start service
	log.Info().Msg("Start account service")

	// load config
	cfgFile, err := config.LoadConfig("./cmd/account/config/config")
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
	conn, err := pgxpool.New(ctx, cfg.Postgres.DnsURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	migrate_db.RunDBMigration(cfg.Postgres.Migration, cfg.Postgres.DnsURL)

	store := db.NewStore(conn)

	// create grpc server
	grpcServer, err := grpc.NewServer(cfg, store)
	if err != nil {
		log.Fatal().Msgf("Create grpc server: %v", err)
	}

	// run grpc server
	go func() {
		if err := grpcServer.Run(); err != nil {
			log.Fatal().Msgf("Run grpc server: %v", err)
		}
	}()

	// create http gateway server
	HTTPGatewayServer, err := http.NewHTTPGatewayServer(cfg, store)
	if err != nil {
		log.Fatal().Msgf("Create http gateway server: %v", err)
	}

	// run http gateway server
	go func() {
		log.Info().Msgf("HTTP Gateway server is running on port %s", cfg.HTTP.Port)
		if err := HTTPGatewayServer.Run(); err != nil {
			log.Fatal().Msgf("Run http gateway server: %v", err)
		}
	}()

	// graceful shutdown
	<-ctx.Done()
	go func() {
		time.Sleep(cfg.GRPC.ShutdownWait * time.Second)
		log.Fatal().Msg("Graceful shutdown timeout")

		HTTPGatewayServer.Shutdown(context.Background())

		grpcServer.GracefulStop()
		doneCh <- struct{}{}
	}()

	<-doneCh
	log.Info().Msg("Account service shutdown")

}
