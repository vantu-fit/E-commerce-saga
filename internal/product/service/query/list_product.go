package query

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
)

type ListProduct struct {
	Offset int64
	Limit  int64
}

type ListProductResult struct {
	ID         uuid.UUID
	CategoryID uuid.UUID
	Name       string
	Price      int64
	Inventory  int64
}

type ListProductHandler QueryHandler[ListProduct, *[]ListProductResult]

type listProductHandler struct {
	store db.Store
}

func NewListProductHandler(store db.Store) ListProductHandler {
	return &listProductHandler{
		store: store,
	}
}

func (h *listProductHandler) Handle(ctx context.Context, cmd ListProduct) (*[]ListProductResult, error) {
	products, err := h.store.ListProducts(ctx, db.ListProductsParams{
		Offset: int32(cmd.Offset),
		Limit:  int32(cmd.Limit),
	})
	if err != nil {
		return nil, err
	}

	result := make([]ListProductResult, len(products))
	for i, product := range products {
		result[i] = ListProductResult{
			ID:         product.ID,
			CategoryID: product.IDCategory,
			Name:       product.Name,
			Price:      product.Price,
			Inventory:  product.Inventory,
		}
	}

	return &result, nil
}
