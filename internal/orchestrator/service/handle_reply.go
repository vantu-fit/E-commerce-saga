package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/vantu-fit/saga-pattern/internal/orchestrator/service/entity"
	"github.com/vantu-fit/saga-pattern/pkg/event"
)

func (s *service) HandleReply(ctx context.Context, msg *kafka.Message) error {
	switch string(msg.Headers[0].Value) {
	case event.UpdateProductInventoryHandler:
		purchaseResult, err := decodePbResponseToEventModel(msg.Value)
		if err != nil {
			return err
		}

		fmt.Println("HandleReply: ", purchaseResult)
		fmt.Println("HandleReply: ", purchaseResult.Succsess)
		if purchaseResult.Succsess {
			err := s.publishPurchaseResult(ctx, &event.PurchaseResult{
				PurchaseID: purchaseResult.Purchase.ID,
				Step:       event.StepUpdateProductInventory,
				Status:     event.StatusSuccess,
			})
			if err != nil {
				return err
			}
			return s.createOrder(ctx, purchaseResult.Purchase)
		}

		err = s.publishPurchaseResult(ctx, &event.PurchaseResult{
			PurchaseID: purchaseResult.Purchase.ID,
			Step:       event.StepUpdateProductInventory,
			Status:     event.StatusFailed,
		})
		if err != nil {
			return err
		}

		return s.rollbackUpdateProductInventory(ctx, purchaseResult.Purchase)
	case event.RollbackProductInventoryHandler:
		purchaseResult, err := decodePbResponseToEventModel(msg.Value)
		if err != nil {
			return err
		}

		if !purchaseResult.Succsess {
			return s.publishPurchaseResult(ctx, &event.PurchaseResult{
				PurchaseID: purchaseResult.Purchase.ID,
				Step:       event.StepUpdateProductInventory,
				Status:     event.StatusRollbackFailed,
			})
		}

		return nil
	case event.CreateOrderHandler:
		purchaseResult, err := decodePbResponseToEventModel(msg.Value)
		if err != nil {
			return err
		}

		if purchaseResult.Succsess {
			err = s.publishPurchaseResult(ctx, &event.PurchaseResult{
				PurchaseID: purchaseResult.Purchase.ID,
				Step:       event.StepCreateOrder,
				Status:     event.StatusSuccess,
			})
			if err != nil {
				return err
			}
			return s.createPayment(ctx, purchaseResult.Purchase)
		}

		err = s.publishPurchaseResult(ctx, &event.PurchaseResult{
			PurchaseID: purchaseResult.Purchase.ID,
			Step:       event.StepCreateOrder,
			Status:     event.StatusFailed,
		})
		if err != nil {
			return err
		}

		return s.rollbackCreateOrder(ctx, purchaseResult.Purchase)
	case event.RollbackOrderHandler:
		purchaseResult, err := decodePbResponseToEventModel(msg.Value)
		if err != nil {
			return err
		}

		if !purchaseResult.Succsess {
			return s.publishPurchaseResult(ctx, &event.PurchaseResult{
				PurchaseID: purchaseResult.Purchase.ID,
				Step:       event.StepCreateOrder,
				Status:     event.StatusRollbackFailed,
			})
		}

		return nil
	case event.CreatePaymentHandler:
		purchaseResult, err := decodePbResponseToEventModel(msg.Value)
		if err != nil {
			return err
		}

		if purchaseResult.Succsess {
			return s.publishPurchaseResult(ctx, &event.PurchaseResult{
				PurchaseID: purchaseResult.Purchase.ID,
				Step:       event.StepCreatePayment,
				Status:     event.StatusSuccess,
			})
		}

		err = s.publishPurchaseResult(ctx, &event.PurchaseResult{
			PurchaseID: purchaseResult.Purchase.ID,
			Step:       event.StepCreatePayment,
			Status:     event.StatusFailed,
		})
		if err != nil {
			return err
		}

		return s.rollbackCreatePayment(ctx, purchaseResult.Purchase)
	case event.RollbackPaymentHandler:
		purchaseResult, err := decodePbResponseToEventModel(msg.Value)
		if err != nil {
			return err
		}

		if !purchaseResult.Succsess {
			return s.publishPurchaseResult(ctx, &event.PurchaseResult{
				PurchaseID: purchaseResult.Purchase.ID,
				Step:       event.StepCreatePayment,
				Status:     event.StatusRollbackFailed,
			})
		}

		return nil

	}

	return nil
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *service) createOrder(ctx context.Context, purchase *entity.Purchase) error {

	pbCreatePurchase := encodeModel2PurchaseRequset(purchase)

	err := s.publishPurchaseResult(ctx, &event.PurchaseResult{
		PurchaseID: purchase.ID,
		Step:       event.StepCreateOrder,
		Status:     event.StatusExecute,
	})
	if err != nil {
		return err
	}

	payload, err := json.Marshal(pbCreatePurchase)
	if err != nil {
		return err
	}

	return s.producer.PublishMessage(ctx, kafka.Message{
		Topic: event.CreateOrderTopic,
		Value: payload,
	})
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *service) rollbackCreateOrder(ctx context.Context, purchase *entity.Purchase) error {
	pbRollbackRequest := encodeModel2PurchaseRequset(purchase)
	payload, err := json.Marshal(pbRollbackRequest)
	if err != nil {
		return err
	}

	err = s.publishPurchaseResult(ctx, &event.PurchaseResult{
		PurchaseID: purchase.ID,
		Step:       event.StepCreateOrder,
		Status:     event.StatusRollback,
	})
	if err != nil {
		return err
	}

	err = s.producer.PublishMessage(ctx, kafka.Message{
		Topic: event.RollbackOrderTopic,
		Value: payload,
	})
	if err != nil {
		log.Error().Msgf("Orchestrator.RollbackFromOrder.RollbackCretaOrder, err: %s", err)
		return err
	}

	return s.rollbackUpdateProductInventory(ctx, purchase)

}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *service) rollbackUpdateProductInventory(ctx context.Context, purchase *entity.Purchase) error {
	pbRollbackRequest := decodePurchase2PbCreatePurchase(purchase)

	err := s.publishPurchaseResult(ctx, &event.PurchaseResult{
		PurchaseID: purchase.ID,
		Step:       event.StepUpdateProductInventory,
		Status:     event.StatusRollback,
	})
	if err != nil {
		return err
	}

	payload, err := json.Marshal(pbRollbackRequest)
	if err != nil {
		return err
	}
	return s.producer.PublishMessage(ctx, kafka.Message{
		Topic: event.RollbackProductInventoryTopic,
		Value: payload,
	})
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *service) createPayment(ctx context.Context, purchase *entity.Purchase) error {
	pbCreatePurchase := encodeModel2PurchaseRequset(purchase)

	err := s.publishPurchaseResult(ctx, &event.PurchaseResult{
		PurchaseID: purchase.ID,
		Step:       event.StepCreatePayment,
		Status:     event.StatusExecute,
	})
	if err != nil {
		return err
	}

	payload, err := json.Marshal(pbCreatePurchase)
	if err != nil {
		return err
	}
	return s.producer.PublishMessage(ctx, kafka.Message{
		Topic: event.CreatePaymentTopic,
		Value: payload,
	})

}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *service) rollbackCreatePayment(ctx context.Context, purchase *entity.Purchase) error {
	pbCreatePurchase := encodeModel2PurchaseRequset(purchase)

	err := s.publishPurchaseResult(ctx, &event.PurchaseResult{
		PurchaseID: purchase.ID,
		Step:       event.StepCreatePayment,
		Status:     event.StatusRollback,
	})
	if err != nil {
		return err
	}

	payload, err := json.Marshal(pbCreatePurchase)
	if err != nil {
		return err
	}

	err = s.producer.PublishMessage(ctx, kafka.Message{
		Topic: event.RollbackPaymentTopic,
		Value: payload,
	})
	if err != nil {
		log.Error().Msgf("Orchestrator.RollbackFromPayment.RollbackCreatePayment, err: %s", err)
		return err
	}

	return s.rollbackCreateOrder(ctx, purchase)
}
