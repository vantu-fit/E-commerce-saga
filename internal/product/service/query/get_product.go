package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
)

type GetProduct struct {
	ProductID uuid.UUID
}

type GetProductResult struct {
	ID          uuid.UUID
	IDAccount   uuid.UUID
	CategoryID  uuid.UUID
	Name        string
	Price       int64
	Inventory   int64
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GetProductHandler QueryHandler[GetProduct, GetProductResult]

type getProductHandler struct {
	store db.Store
}

func NewGetProductHandler(store db.Store) GetProductHandler {
	return &getProductHandler{
		store: store,
	}
}

func (h *getProductHandler) Handle(ctx context.Context, cmd GetProduct) (GetProductResult, error) {
	product, err := h.store.GetProductByID(ctx, cmd.ProductID)
	if err != nil {
		return GetProductResult{}, err
	}

	return GetProductResult{
		ID:          product.ID,
		IDAccount:   product.IDAccount,
		CategoryID:  product.IDCategory,
		Name:        product.Name,
		Price:       product.Price,
		Inventory:   product.Inventory,
		Description: product.Description,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}
