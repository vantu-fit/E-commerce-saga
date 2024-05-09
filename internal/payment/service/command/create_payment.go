package command

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/payment/db/sqlc"
)

type CreatePayment struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Currency   string
	Amount     int64
}

type CreatePaymentHandler CommandHanlder[CreatePayment]

type createPaymentHandler struct {
	store db.Store
}

func NewCreatePaymentHandler(store db.Store) CreatePaymentHandler {
	return &createPaymentHandler{
		store: store,
	}
}

func (h *createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) error {
	_, err := h.store.CreatePayment(ctx, db.CreatePaymentParams{
		ID:         cmd.ID,
		CustomerID: cmd.CustomerID,
		Currency:   cmd.Currency,
		Amount:     cmd.Amount,
	})
	if err != nil {
		return err
	}
	return nil
}
