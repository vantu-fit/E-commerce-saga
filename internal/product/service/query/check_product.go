package query

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
)

type CheckProduct struct {
	ProductID uuid.UUID
	Quantity  int64
}

type CheckProductResult struct {
	ID     uuid.UUID
	Price  int64
	Status bool
}

type CheckProductHandler QueryHandler[CheckProduct, CheckProductResult]

type checkProductHandler struct {
	store db.Store
}

func NewCheckProductHandler(store db.Store) CheckProductHandler {
	return &checkProductHandler{
		store: store,
	}
}

func (h *checkProductHandler) Handle(ctx context.Context, cmd CheckProduct) (CheckProductResult, error) {
	product, err := h.store.GetProductByID(ctx, cmd.ProductID)
	if err != nil {
		return CheckProductResult{}, err
	}

	if product.Inventory < cmd.Quantity {
		return CheckProductResult{
			ID:     product.ID,
			Price:  product.Price,
			Status: false,
		}, nil
	}
	
	return CheckProductResult{
		ID:     product.ID,
		Price:  product.Price,
		Status: true,
	}, nil
}
