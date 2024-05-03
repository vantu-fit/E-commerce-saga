package command

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/order/db/sqlc"
)

type CreateOrderHandler CommandHanlder[CreateOrder]

type CreateOrder struct {
	OrderID    uuid.UUID
	CustomerID uuid.UUID
	Products   *[]PurchasedProduct
}

type PurchasedProduct struct {
	ProductID uuid.UUID
	Quantity  uint64
}

type createOrderHandler struct {
	store db.Store
}

func NewCreateOrderHandler(store db.Store) CreateOrderHandler {
	return &createOrderHandler{
		store: store,
	}
}

func (h *createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) error {
	for _, product := range *cmd.Products {
		arg := db.CreateOrderParams{
			ID:         cmd.OrderID,
			ProductID:  product.ProductID,
			Quantity:   int32(product.Quantity),
			CustomerID: cmd.CustomerID,
		}
		_ , err := h.store.CreateOrder(ctx, arg)
		if err != nil {
			return err
		}
	}
	return nil
}
