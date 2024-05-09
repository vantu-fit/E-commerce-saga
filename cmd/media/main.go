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
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vantu-fit/saga-pattern/cmd/media/config"
	db "github.com/vantu-fit/saga-pattern/internal/media/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/media/grpc"
	"github.com/vantu-fit/saga-pattern/internal/media/http"
	"github.com/vantu-fit/saga-pattern/internal/media/media"
	"github.com/vantu-fit/saga-pattern/internal/media/service"
	grpcclient "github.com/vantu-fit/saga-pattern/pkg/grpc_client"
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
	log.Info().Msg("Start media service")

	// load config
	cfgFile, err := config.LoadConfig("./cmd/media/config/config")
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

	// create minio client
	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.Username, cfg.Minio.Password, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal().Msgf("Create minio client: %v", err)
	}
	minioService := media.New(cfg, minioClient)
	mediaService := service.NewService(store, minioService)

	// create grpc client
	grpcClient := grpcclient.NewClient()

	// create grpc server
	grpcServer := grpc.NewServer(cfg, store, mediaService, grpcClient)

	// run grpc server
	go func() {
		log.Info().Msgf("GRPC server is running on port %s", cfg.GRPC.Port)
		if err := grpcServer.Run(); err != nil {
			log.Fatal().Msgf("Run grpc server: %v", err)
		}
	}()

	// create http gateway server
	httpServer, err := http.NewHTTPGatewayServer(cfg, store, mediaService, grpcClient)
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

	// graceful shutdown
	<-ctx.Done()
	go func() {
		time.Sleep(0 * time.Second)
		log.Fatal().Msg("Graceful shutdown timeout")

		doneCh <- struct{}{}
	}()

	log.Info().Msg("product service shutdown")
	<-doneCh

}
