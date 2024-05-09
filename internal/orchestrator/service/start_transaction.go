package service

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/vantu-fit/saga-pattern/internal/orchestrator/service/entity"
	"github.com/vantu-fit/saga-pattern/pb"
	"github.com/vantu-fit/saga-pattern/pkg/event"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *service) StartTransaction(ctx context.Context, purchase *entity.Purchase) error {
	// init create purchase transaction with step update inventory product
	err := s.publishPurchaseResult(ctx, &event.PurchaseResult{
		PurchaseID: purchase.ID,
		Status:     event.StatusExecute,
		Step:       event.StepUpdateProductInventory,
	})
	if err != nil {
		return err
	}

	// update inventory from product service
	pbCreatePurchase := encodeModel2PurchaseRequset(purchase)

	payload, err := json.Marshal(pbCreatePurchase)
	if err != nil {
		return err
	}

	return s.producer.PublishMessage(ctx, kafka.Message{
		Topic: event.UpdateProductInventoryTopic,
		Value: payload,
	})
}

func (s *service) publishPurchaseResult(ctx context.Context, result *event.PurchaseResult) error {
	pbResult := &pb.PurchaseResult{
		PurchaseId: result.PurchaseID.String(),
		Status:     result.Status,
		Step:       result.Step,
		CreatedAt:  timestamppb.Now(),
	}

	payload, err := json.Marshal(pbResult)
	if err != nil {
		return err
	}

	return s.producer.PublishMessage(ctx, kafka.Message{
		Topic: event.PurchaseResultTopic,
		Value: payload,
	})
}
