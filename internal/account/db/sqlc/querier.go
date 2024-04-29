// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	GetAccountByEmail(ctx context.Context, email string) (Account, error)
	GetSessionById(ctx context.Context, id uuid.UUID) (Session, error)
}

var _ Querier = (*Queries)(nil)