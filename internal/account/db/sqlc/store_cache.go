package db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/vantu-fit/saga-pattern/cmd/account/config"
	"github.com/vantu-fit/saga-pattern/pkg/cache"
)

type StoreCache struct {
	Store
	lc     cache.LocalCache
	rc     cache.RedisCache
	config *config.Config
}

const (
	cuckooFilter = "account_cuckoo_filter"
	dummnyItem   = "dummy_item"
	mutexKey     = "mutex:"

	getaccountKey = "account:"
	sessionKey    = "session:"
)

func NewStoreCache(
	store Store,
	lc cache.LocalCache,
	rc cache.RedisCache,
	config *config.Config,
) Store {
	exist, err := rc.CFExist(context.Background(), cuckooFilter, dummnyItem)
	if err != nil {
		log.Error().Msgf("Account: failed to check cuckoo filter existence, err: %s", err)
		return store
	}

	if !exist {
		err = rc.CFReserve(context.Background(), cuckooFilter, 1000, 4, 1000)
		if err != nil {
			log.Error().Msgf("Account: failed to reserve cuckoo filter, err: %s", err)
			return store
		}

		err = rc.CFAdd(context.Background(), cuckooFilter, dummnyItem)
		if err != nil {
			log.Error().Msgf("Account: failed to add dummy item to cuckoo filter, err: %s", err)
			return store
		}

		log.Info().Msg("Account: reserved cuckoo filter")
	}

	return StoreCache{
		Store:  store,
		lc:     lc,
		rc:     rc,
		config: config,
	}
}

func (sc StoreCache) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	account, err := sc.Store.CreateAccount(ctx, arg)
	if err != nil {
		log.Error().Msgf("Account: failed to create account, err: %s", err)
		return account, err
	}
	err = sc.rc.CFAdd(ctx, cuckooFilter, account.ID.String())
	if err != nil {
		log.Error().Msgf("Account: failed to add account to cuckoo filter, err: %s", err)
		return account, err
	}
	err = sc.lc.Set(getaccountKey+account.ID.String(), &account)
	if err != nil {
		log.Error().Msgf("Account: failed to set account to local cache, err: %s", err)
	}
	return account, err
}

func (sc StoreCache) GetAccount(ctx context.Context, id uuid.UUID) (Account, error) {
	var account Account
	key := getaccountKey + id.String()
	// Get from local cache
	ok, err := sc.lc.Get(key, &account)
	if ok && err == nil {
		return account, err
	}
	// Check if account exist in cuckoo filter
	exist, err := sc.rc.CFExist(ctx, cuckooFilter, id.String())
	if !exist && err == nil {
		account, err := sc.Store.GetAccount(ctx, id)
		if err != nil {
			return account, err
		}
		err = sc.rc.CFAdd(ctx, cuckooFilter, id.String())
		if err != nil {
			log.Error().Msgf("Account: failed to add account to cuckoo filter, err: %s", err)
		}

		err = sc.lc.Set(key, &account)
		if err != nil {
			log.Error().Msgf("Account: failed to set account to local cache, err: %s", err)
		}

		err = sc.rc.Set(ctx, key, &account)
		if err != nil {
			log.Error().Msgf("Account: failed to set account to redis cache, err: %s", err)
		}

		return account, nil
	}
	// Get from redis cache
	ok, err = sc.rc.Get(ctx, key, &account)
	if ok && err == nil {
		err = sc.lc.Set(key, &account)
		if err != nil {
			log.Error().Msgf("Account: failed to set account to local cache, err: %s", err)
		}
		return account, nil
	}
	// Lock mutex to prevent cache stampede
	mu := sc.rc.GetMutex(mutexKey + key)
	err = mu.Lock()
	if err != nil {
		return account, err
	}
	defer mu.Unlock()

	// Get again to prevent new update
	ok, err = sc.rc.Get(ctx, key, &account)
	if ok && err == nil {
		err = sc.lc.Set(key, &account)
		if err != nil {
			log.Error().Msgf("Account: failed to set account to local cache, err: %s", err)
		}
		return account, nil
	}
	// Get from database
	account, err = sc.Store.GetAccount(ctx, id)
	if err != nil {
		return account, err
	}
	// Set to redis cache
	err = sc.rc.Set(ctx, key, &account)
	if err != nil {
		log.Error().Msgf("Account: failed to set account to local cache, err: %s", err)
	}
	// Set to local cache
	err = sc.lc.Set(key, &account)
	if err != nil {
		log.Error().Msgf("Account: failed to set account to redis cache, err: %s", err)
	}

	return account, nil
}

func (sc StoreCache) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	session := &Session{
		ID:           arg.ID,
		UserID:       arg.UserID,
		RefreshToken: arg.RefreshToken,
		UserAgent:    arg.UserAgent,
		ClientIp:     arg.ClientIp,
	}
	key := sessionKey + session.ID.String()

	err := sc.rc.Set(context.Background(), key, session, int(sc.config.PasetoConfig.RefreshTokenExpire))
	if err != nil {
		log.Error().Msgf("Account: failed to set session, err: %s", err)
		return Session{}, err
	}

	return *session, nil
}

func (sc StoreCache) GetSessionById(ctx context.Context, id uuid.UUID) (Session, error) {
	var session = &Session{}
	key := sessionKey + id.String()

	ok, err := sc.lc.Get(key, session)
	if ok && err == nil {
		return *session, nil
	}

	ok, err = sc.rc.Get(ctx, key, session)
	if ok && err == nil {
		err := sc.lc.Set(key, session)
		if err != nil {
			log.Error().Msgf("Account: failed to set session to local cache, err: %s", err)
		}
		return *session, nil
	}
	log.Info().Msgf("Account: Unauthorizate : %s", key)

	return *session, errors.New("Unauthorizate")
}
