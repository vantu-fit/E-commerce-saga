package db_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	db "github.com/vantu-fit/saga-pattern/internal/account/db/sqlc"
)

func TestCreateAccount(t *testing.T) {
	account := createRandomAccount()
	arg := db.CreateAccountParams{
		FirstName:   account.FirstName,
		LastName:    account.LastName,
		Email:       account.Email,
		Password:    account.Password,
		Address:     account.Address,
		PhoneNumber: account.PhoneNumber,
	}
	accountRes, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accountRes)
	require.Equal(t, account.FirstName, accountRes.FirstName)
	require.Equal(t, account.LastName, accountRes.LastName)
	require.Equal(t, account.Email, accountRes.Email)
	require.Equal(t, account.Address, accountRes.Address)
	require.Equal(t, account.PhoneNumber, accountRes.PhoneNumber)
	require.NotZero(t, accountRes.ID)
	require.NotZero(t, accountRes.CreatedAt)
	require.NotZero(t, accountRes.UpdatedAt)

}

func TestCreateSession(t *testing.T) {
	account := createRandomAccount()
	arg := db.CreateAccountParams{
		FirstName:   account.FirstName,
		LastName:    account.LastName,
		Email:       account.Email,
		Password:    account.Password,
		Address:     account.Address,
		PhoneNumber: account.PhoneNumber,
	}
	accountRes, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accountRes)

	argSession := db.CreateSessionParams{
		ID:           uuid.New(),
		Email:        accountRes.Email,
		RefreshToken: RandomString(32),
		UserAgent:    "",
		ClientIp:     "",
		ExpiresAt:    time.Now().Add(time.Minute * 15),
	}
	session, err := testStore.CreateSession(context.Background(), argSession)
	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, argSession.ID, session.ID)
	require.Equal(t, argSession.Email, session.Email)
	require.Equal(t, argSession.RefreshToken, session.RefreshToken)
	require.Equal(t, argSession.UserAgent, session.UserAgent)
	require.Equal(t, argSession.ClientIp, session.ClientIp)
	require.NotZero(t, session.CreatedAt)
}

func createRandomAccount() *db.Account {
	return &db.Account{
		Email:       RandomEmail(),
		FirstName:   RandomString(8),
		LastName:    RandomString(8),
		Password:    RandomString(8),
		Address:     RandomString(8),
		PhoneNumber: RandomString(8),
	}
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func RandomEmail() string {
	return RandomString(8) + "@gmail.com"
}
