package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/rs/zerolog/log"
)

const (
	DefaultExpiration  = 600
	CleanWindowMinutes = 3
)

type LocalCache interface {
	Get(key string, value interface{}) (bool, error)
	Set(key string, value interface{}) error
	Delete(key string) error
}

type localCache struct {
	cache *bigcache.BigCache
}

func NewLocalCache(ctx context.Context, expirationTime uint64) (LocalCache, error) {
	if expirationTime == 0 {
		expirationTime = DefaultExpiration
	}

	config := bigcache.Config{
		LifeWindow:  time.Duration(expirationTime),
		CleanWindow: time.Minute * CleanWindowMinutes,

		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.

		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,

		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 500,

		// prints information about additional memory allocation
		Verbose: true,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 8192,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		OnRemoveWithReason: nil,
		
	}

	cache, err := bigcache.New(ctx, config)
	if err != nil {
		log.Error().Msgf("Create local cache: %v", err)
		return nil, err
	}

	return &localCache{
		cache: cache,
	}, nil
}

func (lc *localCache) Get(key string, value interface{}) (bool, error) {
	val, err := lc.cache.Get(key)
	if err == bigcache.ErrEntryNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if err = json.Unmarshal(val, value); err != nil {
		return false, err
	}

	log.Info().Msgf("Get key local cache: %s, value: %v", key, value)

	return true, nil
}

func (lc *localCache) Set(key string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	log.Info().Msgf("Set key local cache: %s, value: %v", key, string(val))
	
	return lc.cache.Set(key, val)
}

func (lc *localCache) Delete(key string) error {
	return lc.cache.Delete(key)
}
