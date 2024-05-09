package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vantu-fit/saga-pattern/internal/purchase/event"
	"github.com/vantu-fit/saga-pattern/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreatePurchase struct {
	*pb.CreatePurchaseRequestApi
	CustomerId uuid.UUID
	Amount     int64
	PurchaseId uuid.UUID
}

type CreatePurchaseHandler CommandHandler[CreatePurchase]

type createPurchaseHanler struct {
	event event.EventHanlder
}

func NewCreatePurchaseHanler(
	event event.EventHanlder,
) CreatePurchaseHandler {
	return &createPurchaseHanler{
		event: event,
	}

}

func (h *createPurchaseHanler) Handle(ctx context.Context, cmd CreatePurchase) error {
	pbCreatePurchase := pb.CreatePurchaseRequest{
		PurchaseId: cmd.PurchaseId.String(),
		Purchase: &pb.Purchase{
			Order: &pb.Order{
				CustomerId: cmd.CustomerId.String(),
				OrderItems: cmd.OrderItems,
			},
			Payment: &pb.Payment{
				CurrencyCode: cmd.Payment.CurrencyCode,
				Amount:       uint64(cmd.Amount),
			},
		},
		Timestamp: timestamppb.Now(),
	}
	fmt.Println("CreatePurchaseHandler")
	return h.event.ProduceCreatePurchaseEvent(ctx, &pbCreatePurchase)

}
