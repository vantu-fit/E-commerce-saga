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
	"github.com/vantu-fit/saga-pattern/cmd/account/config"
	db "github.com/vantu-fit/saga-pattern/internal/account/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/account/grpc"
	"github.com/vantu-fit/saga-pattern/internal/account/http"

	"github.com/vantu-fit/saga-pattern/pkg/cache"
	migrate_db "github.com/vantu-fit/saga-pattern/pkg/migrate"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	kafkaClient "github.com/vantu-fit/saga-pattern/pkg/kafka"
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

	// create redis cluster client
	log.Info().Msgf("Redis cluster is running on %v", cfg.RedisCache.Address)
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:         cfg.RedisCache.Address,
		Password:      cfg.RedisCache.Password,
		PoolSize:      cfg.RedisCache.PoolSize,
		MaxRetries:    cfg.RedisCache.MaxRetries,
		RouteByLatency: true,
		ReadOnly:      true,
		RouteRandomly: true,
	})
	
	// ping
	err = redisClient.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		log.Fatal().Msgf("Redis cluster ping: %v", err)
	}
	// create local cache
	localCache , err := cache.NewLocalCache(ctx, 0)
	if err != nil {
		log.Fatal().Msgf("Create local cache: %v", err)
	}
	// create redis cache
	redisCache := cache.NewRedisCache(redisClient , time.Duration(cfg.RedisCache.ExpirationTime) * time.Second)
	// create store cache
	storeCache := db.NewStoreCache(store , localCache , redisCache , cfg)
	
	// create kafka producer
	producer := kafkaClient.NewProducer(cfg.Kafka.Brokers)

	// create grpc server
	grpcServer, err := grpc.NewServer(cfg, storeCache , producer)
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
	HTTPGatewayServer, err := http.NewHTTPGatewayServer(cfg, storeCache , producer)
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
