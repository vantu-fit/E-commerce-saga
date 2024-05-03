package query

import (
	"context"
	"errors"

	"github.com/google/uuid"
	db "github.com/vantu-fit/saga-pattern/internal/order/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/order/service/object"
	"github.com/vantu-fit/saga-pattern/pb"
	grpcclient "github.com/vantu-fit/saga-pattern/pkg/grpc_client"
)

type GetOrderDetailHandler QueryHanler[GetOrderDetail, *object.DetailedOrder]

type GetOrderDetail struct {
	OrderID uuid.UUID
}

type getOrderDetailHandler struct {
	store      db.Store
	grpcClient *grpcclient.Client
}

func NewGetOrderDetailHandler(store db.Store, grpcClient *grpcclient.Client) GetOrderDetailHandler {
	return &getOrderDetailHandler{
		store:      store,
		grpcClient: grpcClient,
	}
}

func (h *getOrderDetailHandler) Handle(ctx context.Context, cmd GetOrderDetail) (*object.DetailedOrder, error) {
	orders, err := h.store.GetOrder(ctx, cmd.OrderID)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("No order id match!")
	}

	purchasedProduct := make([]object.DetailedPurchasedProduct , len(orders))

	detailOrder := &object.DetailedOrder{
		ID:         cmd.OrderID,
		CustomerID: orders[0].CustomerID,
		PurchasedProducts: &purchasedProduct,
	}

	for i , order := range orders {
		product, err := h.grpcClient.ProductClient.GetProductByID(ctx, &pb.GetProductByIDRequest{
			Id: order.ProductID.String(),
		})
		if err != nil {
			return nil, err
		}
		purchasedProduct[i] = object.DetailedPurchasedProduct{
			ID: uuid.MustParse(product.Product.Id),
			CategoryID: uuid.MustParse(product.Product.CategoryId),
			Name: product.Product.Name,
			BrandName: product.Product.BrandName,
			Description: product.Product.Description,
			Price: uint32(product.Product.Price),
			Quantity: uint32(product.Product.Inventory),
		}
	}

	return detailOrder, nil
}
