// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error)
	CreateIdempotency(ctx context.Context, arg CreateIdempotencyParams) (Idempotency, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	GetIdempotencyKey(ctx context.Context, id uuid.UUID) ([]Idempotency, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (Product, error)
	GetProductCategory(ctx context.Context, name string) ([]GetProductCategoryRow, error)
	GetProductForUpdate(ctx context.Context, id uuid.UUID) (Product, error)
	UpadateProduct(ctx context.Context, arg UpadateProductParams) (Product, error)
	UpdateIdempotency(ctx context.Context, arg UpdateIdempotencyParams) ([]Idempotency, error)
}

var _ Querier = (*Queries)(nil)
