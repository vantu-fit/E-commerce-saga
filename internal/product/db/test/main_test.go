package db_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vantu-fit/saga-pattern/cmd/product/config"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pkg/cache"
)

var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

var testStore db.Store

func TestMain(m *testing.M) {

	// start service
	log.Info().Msg("Start product service")

	// load config
	cfgFile, err := config.LoadConfig("../../../../cmd/product/config/config")
	if err != nil {
		log.Fatal().Msgf("Load config: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatal().Msgf("Parse config: %v", err)
	}

	// create context for gracefull shutdown
	ctx, stop := signal.NotifyContext(context.Background(), interuptSignals...)
	defer stop()

	// run migrate DB
	poolConfig, err := pgxpool.ParseConfig(cfg.Postgres.DnsURL)
	if err != nil {
		log.Fatal().Msgf("Parse pgx pool config: %v", err)
	}
	poolConfig.MaxConns = 500 // Set maximum connections in the pool
	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

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

	// create store cache
	testStore = db.NewStoreCache(store, localCache, redisCache)

	m.Run()
}
