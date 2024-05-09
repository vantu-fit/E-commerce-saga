package service

import (
	db "github.com/vantu-fit/saga-pattern/internal/payment/db/sqlc"
	"github.com/vantu-fit/saga-pattern/internal/payment/service/command"
	"github.com/vantu-fit/saga-pattern/internal/payment/service/query"
)

type Command struct {
	CreatePayment command.CreatePaymentHandler
	DeletePayment command.DeletePaymentHandler
}

type Query struct {
	GetPayment query.GetPaymentHandler
}

type Service struct {
	Command
	Query
}

func NewService(store db.Store) *Service {
	return &Service{
		Command: Command{
			CreatePayment: command.NewCreatePaymentHandler(store),
			DeletePayment: command.NewDeletePaymentHandler(store),
		},
		Query: Query{
			GetPayment: query.NewGetPaymentHandler(store),
		},
	}
}
