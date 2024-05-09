package service

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/vantu-fit/saga-pattern/internal/orchestrator/service/entity"
	"github.com/vantu-fit/saga-pattern/pb"
	"github.com/vantu-fit/saga-pattern/pkg/event"
)


// func getPbPurchaseStep(step string) pb.PurchaseStep {
// 	switch step {
// 	case event.StepUpdateProductInventory:
// 		return pb.PurchaseStep_UPDATE_PRODUCT_INVENTORY
// 	case event.StepCreateOrder:
// 		return pb.PurchaseStep_CREATE_ORDER
// 	case event.StepCreatePayment:
// 		return pb.PurchaseStep_CREATE_PAYMENT
// 	}
// 	return -1
// }

// func getPbPurchaseStatus(status string) pb.PurchaseStatus {
// 	switch status {
// 	case event.StatusExecute:
// 		return pb.PurchaseStatus_EXECUTE
// 	case event.StatusSuccess:
// 		return pb.PurchaseStatus_SUCCESS
// 	case event.StatusFailed:
// 		return pb.PurchaseStatus_FAILED
// 	case event.StatusRollback:
// 		return pb.PurchaseStatus_ROLLBACK
// 	case event.StatusRollbackFailed:
// 		return pb.PurchaseStatus_ROLLBACK_FAILED
// 	}
// 	return -1
// }

func encodeModel2PurchaseRequset(purchase *entity.Purchase) *pb.CreatePurchaseRequest {
	orderItems := make([]*pb.PurchaseOrderItem, len(*purchase.Order.OrderItems))
	for i, orderItem := range *purchase.Order.OrderItems {
		orderItems[i] = &pb.PurchaseOrderItem{
			ProductId: orderItem.ID.String(),
			Quantity:  orderItem.Quantity,
		}
	}

	pbCreatePurchase := &pb.CreatePurchaseRequest{
		PurchaseId: purchase.ID.String(),
		Purchase: &pb.Purchase{
			Order: &pb.Order{
				CustomerId: purchase.Order.CustomerID.String(),
				OrderItems: orderItems,
			},
			Payment: &pb.Payment{
				CurrencyCode: purchase.Payment.CurrentcyCode,
				Amount:       purchase.Payment.Amount,
			},
		},
	}

	return pbCreatePurchase
}

func decodePbResponseToEventModel(data []byte) (*event.CreatePurchaseResponse, error) {
	var pbResult pb.CreatePurchaseResponse
	err := json.Unmarshal(data, &pbResult)
	if err != nil {
		return nil, err
	}

	orderItems := make([]entity.OrderItem, len(pbResult.Purchase.Order.OrderItems))
	for i, item := range pbResult.Purchase.Order.OrderItems {
		orderItems[i] = entity.OrderItem{
			ID:       uuid.MustParse(item.ProductId),
			Quantity: item.Quantity,
		}
	}

	return &event.CreatePurchaseResponse{
		Purchase: &entity.Purchase{
			ID: uuid.MustParse(pbResult.PurchaseId),
			Order: &entity.Order{
				ID:         uuid.MustParse(pbResult.PurchaseId),
				CustomerID: uuid.MustParse(pbResult.Purchase.Order.CustomerId),
				OrderItems: &orderItems,
			},
			Payment: &entity.Payment{
				ID:            uuid.MustParse(pbResult.PurchaseId),
				CurrentcyCode: pbResult.Purchase.Payment.CurrencyCode,
				Amount:        pbResult.Purchase.Payment.Amount,
			},
		},
		Succsess: pbResult.Success,
		Error:   pbResult.ErrorMessage,
	}, nil
}

func decodePurchase2PbCreatePurchase(purchase *entity.Purchase) *pb.CreatePurchaseRequest {
	orderItems := make([]*pb.PurchaseOrderItem, len(*purchase.Order.OrderItems))
	for i, orderItem := range *purchase.Order.OrderItems {
		orderItems[i] = &pb.PurchaseOrderItem{
			ProductId: orderItem.ID.String(),
			Quantity:  orderItem.Quantity,
		}
	}

	pbCreatePurchase := &pb.CreatePurchaseRequest{
		PurchaseId: purchase.ID.String(),
		Purchase: &pb.Purchase{
			Order: &pb.Order{
				CustomerId: purchase.Order.CustomerID.String(),
				OrderItems: orderItems,
			},
			Payment: &pb.Payment{
				CurrencyCode: purchase.Payment.CurrentcyCode,
				Amount:       purchase.Payment.Amount,
			},
		},
	}

	return pbCreatePurchase
}
