package command

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/payment/db/sqlc"
)

type DeletePayment struct {
	ID uuid.UUID
}

type DeletePaymentHandler CommandHanlder[DeletePayment]

type deletePaymentHandler struct {
	store db.Store
}

func NewDeletePaymentHandler(store db.Store) DeletePaymentHandler {
	return &deletePaymentHandler{
		store: store,
	}
}

func (h *deletePaymentHandler) Handle(ctx context.Context, cmd DeletePayment) error {
	_, err := h.store.DeletePayment(ctx, cmd.ID )
	
	if err != nil {
		return err
	}
	return nil
}
