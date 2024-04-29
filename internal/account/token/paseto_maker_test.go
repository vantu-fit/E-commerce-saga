package token

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	symmetricKey := randString(32)
	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)
	require.NotNil(t, maker)

	payload, err := NewPayload(uuid.New(), "vantu", time.Minute*5)
	require.NoError(t, err)

	token, tokenPayload, err := maker.CreateToken(payload.ID , payload.Email, payload.ExpiredAt.Sub(payload.IssuedAt))
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, tokenPayload)
	require.Equal(t, payload.Email, tokenPayload.Email)
	require.Equal(t, payload.IssuedAt.Unix(), tokenPayload.IssuedAt.Unix())
	require.Equal(t, payload.ExpiredAt.Unix(), tokenPayload.ExpiredAt.Unix())
	require.NotEmpty(t, tokenPayload.ID)
}

func TestVerifyToken(t *testing.T) {
	symmetricKey := randString(32)
	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)
	require.NotNil(t, maker)

	payload, err := NewPayload(uuid.New(), "vantu", time.Minute*5)

	token, newPayload, err := maker.CreateToken(payload.ID , payload.Email, payload.ExpiredAt.Sub(payload.IssuedAt))
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, newPayload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, newPayload.Email, payload.Email)
	require.Equal(t, newPayload.IssuedAt.Unix(), payload.IssuedAt.Unix())
	require.Equal(t, newPayload.ExpiredAt.Unix(), payload.ExpiredAt.Unix())
	require.NotEmpty(t, payload.ID)
}

func TestInvalidToken(t *testing.T) {
	symmetricKey := randString(32)
	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)
	require.NotNil(t, maker)

	payload, err := NewPayload(uuid.New(), "vantu", time.Minute*5)

	token, _, err := maker.CreateToken(payload.ID , payload.Email, payload.ExpiredAt.Sub(payload.IssuedAt))
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Change token
	runes := []rune(token)
	runes[rand.Intn(len(runes))]++
	invalidToken := string(runes)

	payload, err = maker.VerifyToken(invalidToken)
	require.Error(t, err)
	require.EqualError(t, err, "invalid token authentication")
	require.Nil(t, payload)
}

func TestExpiredToken(t *testing.T) {
	symmetricKey := randString(32)
	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)
	require.NotNil(t, maker)

	payload, err := NewPayload(uuid.New(), "vantu", -time.Minute*5)

	token, _, err := maker.CreateToken(payload.ID , payload.Email, payload.ExpiredAt.Sub(payload.IssuedAt))
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, "token has expried")
	require.Nil(t, payload)
}

func randString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
