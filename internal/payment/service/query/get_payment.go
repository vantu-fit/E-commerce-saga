package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/payment/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/payment/service/entity"
)

type GetPayment struct {
	ID uuid.UUID
}

type GetPaymentResponse struct {
	Payment   *entity.Payment
	CreatedAt time.Time
}

type GetPaymentHandler QueryHandler[GetPayment, *GetPaymentResponse]

type getPaymentHandler struct {
	store db.Store
}

func NewGetPaymentHandler(store db.Store) GetPaymentHandler {
	return &getPaymentHandler{
		store: store,
	}
}

func (h *getPaymentHandler) Handle(ctx context.Context, cmd GetPayment) (*GetPaymentResponse, error) {
	payment, err := h.store.GetPaymentById(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	response := GetPaymentResponse{
		Payment: &entity.Payment{
			ID:            payment.ID,
			CurrentcyCode: payment.Currency,
			Amount:        uint64(payment.Amount),
		},
		CreatedAt: payment.CreatedAt,
	}

	return &response, nil
}
