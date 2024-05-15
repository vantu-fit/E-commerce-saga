package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisCache interface {
	Get(ctx context.Context, key string, value interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{} , expiration ...int) error
	Delete(ctx context.Context, key string) error
	CFReserve(ctx context.Context, key string, capacity, bucketSize, maxIterations int64) error
	CFAdd(ctx context.Context, key string, value interface{}) error
	CFExist(ctx context.Context, key string, value interface{}) (bool, error)
	CFDel(ctx context.Context, key string, value interface{}) error
	GetMutex(muxtexName string) *redsync.Mutex
	ExecIncrbyXPipeline(ctx context.Context, payloads *[]RedisIncrbyXPayload) error
}

type redisCache struct {
	cache      *redis.ClusterClient
	rs         *redsync.Redsync
	expiration time.Duration
}

type RedisIncrbyXPayload struct {
	Key   string
	Value int64
}

func NewRedisCache(cache *redis.ClusterClient, expirationTime time.Duration) (RedisCache) {
	pool := goredis.NewPool(cache)
	rs := redsync.New(pool)
	return &redisCache{
		cache:      cache,
		rs:         rs,
		expiration: expirationTime,
	}
}

func (rc *redisCache) Get(ctx context.Context, key string, value interface{}) (bool, error) {
	val, err := rc.cache.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if err = json.Unmarshal([]byte(val), value); err != nil {
		return false, err
	}

	log.Info().Msgf("Get key redis cache: %s, value: %v", key, value)

	return true, nil
}

func (rc *redisCache) Set(ctx context.Context, key string, value interface{} , expiration ...int) error {
	expirationTime := rc.expiration
	if len(expiration) > 0 {
		expirationTime =time.Duration(expiration[0]) *60 *time.Second
	}
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	log.Info().Msgf("Set key redis cache: %s, value: %v", key, string(val))
	return rc.cache.Set(ctx, key, val , expirationTime).Err()
}

func (rc *redisCache) Delete(ctx context.Context, key string) error {
	return rc.cache.Del(ctx, key).Err()
}

func (rc *redisCache) CFReserve(ctx context.Context, key string, capacity, bucketSize, maxIterations int64) error {
	return rc.cache.CFReserveWithArgs(ctx, key, &redis.CFReserveOptions{
		Capacity:      capacity,
		BucketSize:    bucketSize,
		MaxIterations: maxIterations,
	}).Err()
}

func (rc *redisCache) CFAdd(ctx context.Context, key string, value interface{}) error {
	return rc.cache.CFAdd(ctx, key, value).Err()
}

func (rc *redisCache) CFExist(ctx context.Context, key string, value interface{}) (bool, error) {
	val, err := rc.cache.CFExists(ctx, key, value).Result()
	if err != nil {
		return false, nil
	}
	return val, nil
}

func (rc *redisCache) CFDel(ctx context.Context, key string, value interface{}) error {
	return rc.cache.CFDel(ctx, key, value).Err()
}

func (rc *redisCache) GetMutex(muxtexName string) *redsync.Mutex {
	return rc.rs.NewMutex(muxtexName)
}

var incrByX = redis.NewScript(`
local exists = redis.call("EXISTS", KEYS[1])
if exists == 1 then 
	return redis.call("INCRBY", KEYS[1], ARGV[1])
else
	return redis.call("SET", KEYS[1], ARGV[1])
end
`)

func (rc *redisCache) ExecIncrbyXPipeline(ctx context.Context, payloads *[]RedisIncrbyXPayload) error {
	pipe := rc.cache.Pipeline()
	executedCmds := make([]*redis.Cmd, len(*payloads))
	for i, payload := range *payloads {
		executedCmds[i] = incrByX.Run(ctx, rc.cache, []string{payload.Key}, payload.Value)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	for _, cmd := range executedCmds {
		if err = cmd.Err(); err != nil {
			return err
		}
	}

	return nil
}
