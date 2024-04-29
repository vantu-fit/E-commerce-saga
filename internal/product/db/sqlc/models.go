// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type Idempotency struct {
	ID         uuid.UUID `json:"id"`
	ProductID  uuid.UUID `json:"product_id"`
	Quantity   int32     `json:"quantity"`
	Rollbacked bool      `json:"rollbacked"`
	CreatedAt  time.Time `json:"created_at"`
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	IDCategory  uuid.UUID `json:"id_category"`
	IDAccount   uuid.UUID `json:"id_account"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BrandName   string    `json:"brand_name"`
	Price       int32     `json:"price"`
	Inventory   int32     `json:"inventory"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}
