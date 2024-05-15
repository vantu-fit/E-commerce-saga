package cache_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	db "github.com/vantu-fit/saga-pattern/internal/account/db/sqlc"
	"github.com/vantu-fit/saga-pattern/pkg/cache"
)

func TestLocalcache(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cache, err := cache.NewLocalCache(context.Background(), 600)
	require.NoError(t, err)

	session := &db.Session{
		ID:           uuid.New(),
		UserID:       uuid.New(),
		RefreshToken: "refresh_token",
		UserAgent:    "",
		ClientIp:     "",
	}

	err = cache.Set("session:"+session.RefreshToken, session)
	require.NoError(t, err)

	ok, err := cache.Get("session:"+session.RefreshToken, session)
	require.NoError(t, err)
	require.True(t, ok)

}
