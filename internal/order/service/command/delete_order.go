package command

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/order/db/sqlc"
)

type DeleteOrderHandler CommandHanlder[DeleteOrder]

type DeleteOrder struct {
	OrderID uuid.UUID
}

type deleteOrderHandler struct {
	store db.Store
}

func NewDeleteOrderHandler(store db.Store) DeleteOrderHandler {
	return &deleteOrderHandler{
		store: store,
	}
}

func (h *deleteOrderHandler) Handle(ctx context.Context, cmd DeleteOrder) error {
	_, err := h.store.DeleteOrder(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	return nil
}
