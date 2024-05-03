package service

import (
	db "github.com/vantu-fit/saga-pattern/internal/order/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/order/service/command"
	"github.com/vantu-fit/saga-pattern/internal/order/service/query"
	grpcclient "github.com/vantu-fit/saga-pattern/pkg/grpc_client"
)

type Service struct {
	Command
	Query
}

type Command struct {
	CreateOrder command.CreateOrderHandler
	DeleteOrder command.DeleteOrderHandler
}

type Query struct {
	GetOrderDetail query.GetOrderDetailHandler
}

func NewOrderService(store db.Store, grpcClient *grpcclient.Client) *Service {
	return &Service{
		Command: Command{
			CreateOrder: command.NewCreateOrderHandler(store),
			DeleteOrder: command.NewDeleteOrderHandler(store),
		},
		Query: Query{
			GetOrderDetail: query.NewGetOrderDetailHandler(store , grpcClient),
		},
	}
}
