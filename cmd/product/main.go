package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vantu-fit/saga-pattern/cmd/product/config"
	"github.com/vantu-fit/saga-pattern/internal/cache"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/product/grpc"
	"github.com/vantu-fit/saga-pattern/internal/product/http"

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

	// start service
	log.Info().Msg("Start product service")

	// load config
	cfgFile, err := config.LoadConfig("./cmd/product/config/config")
	if err != nil {
		log.Fatal().Msgf("Load config: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatal().Msgf("Parse config: %v", err)
	}

	log.Debug().Msgf("redis addr: %s",  cfg.RedisCache.Address[0])

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

	// create local cache
	localCache, err := cache.NewLocalCache(ctx, cfg.LocalCache.ExpirationTime)
	if err != nil {
		log.Fatal().Msgf("Create local cache: %v", err)
	}

	// create redis cache
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:         cfg.RedisCache.Address,
		Password:      cfg.RedisCache.Password,
		PoolSize:      cfg.RedisCache.PoolSize,
		MaxRetries:    cfg.RedisCache.MaxRetries,
		ReadOnly:      true,
		RouteRandomly: true,
	})

	// check redis connection
	err = redisClient.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		log.Fatal().Msgf("Redis ping: %v", err)
	}

	redisCache := cache.NewRedisCache(redisClient, time.Duration(cfg.RedisCache.ExpirationTime)*time.Second)

	_ = localCache
	_ = redisCache

	// create grpc client
	grpcClient := grpc.NewClient(cfg)

	// run grpc client
	go func() {
		if err := grpcClient.RunAccountClient(doneCh); err != nil {
			log.Fatal().Msgf("Run grpc client: %v", err)
		}
	}()

	// create grpc server
	grpcServer, err := grpc.NewServer(cfg, store, grpcClient)
	if err != nil {
		log.Fatal().Msgf("Create grpc server: %v", err)
	}

	// run grpc server
	go func() {
		log.Info().Msgf("GRPC server is running on port %s", cfg.GRPC.Port)
		if err := grpcServer.Run(); err != nil {
			log.Fatal().Msgf("Run grpc server: %v", err)
		}
	}()

	// create http gateway server
	HTTPGatewayServer, err := http.NewHTTPGatewayServer(cfg, store, grpcClient)
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

	log.Info().Msg("product service shutdown")
	<-doneCh

}
