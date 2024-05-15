package cache_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/vantu-fit/saga-pattern/cmd/account/config"
	db "github.com/vantu-fit/saga-pattern/internal/account/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pkg/cache"
)

func TestGet(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// load config
	cfgFile, err := config.LoadConfig("../../cmd/account/config/config")
	if err != nil {
		log.Fatal().Msgf("Load config: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatal().Msgf("Parse config: %v", err)
	}

	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          cfg.RedisCache.Address,
		Password:       cfg.RedisCache.Password,
		PoolSize:       cfg.RedisCache.PoolSize,
		MaxRetries:     cfg.RedisCache.MaxRetries,
		RouteByLatency: true,
		ReadOnly:       true,
		RouteRandomly:  true,
	})

	// create redis cache
	cacheRedis := cache.NewRedisCache(redisClient, time.Duration(cfg.RedisCache.ExpirationTime)*time.Second)

	// set key
	session := db.Session{
		ID:           uuid.New(),
		UserID:       uuid.New(),
		RefreshToken: uuid.New().String(),
		UserAgent:    "",
		ClientIp:     "",
	}

	err = cacheRedis.Set(context.Background(), "session:"+session.ID.String(), session, int(cfg.RedisCache.ExpirationTime))
	require.NoError(t, err)

	// get key
	var val = &db.Session{}
	ok, err := cacheRedis.Get(context.Background(), "session:"+session.ID.String(), val)
	require.NoError(t, err)
	require.True(t, ok)
}
