package service

import (
	"github.com/vantu-fit/saga-pattern/internal/purchase/event"
	"github.com/vantu-fit/saga-pattern/internal/purchase/service/command"
	"github.com/vantu-fit/saga-pattern/internal/purchase/service/query"
)

type Command struct {
	CreatePurchase command.CreatePurchaseHandler
}

type Query struct {
	CheckProducts query.CheckProductHandler
}

type Service struct {
	Command
	Query
}

func NewService(
	event event.EventHanlder,
) *Service {
	return &Service{
		Command: Command{
			CreatePurchase: command.NewCreatePurchaseHanler(event),
		},
		Query: Query{
			CheckProducts: query.NewCheckProductHandler(),
		},
	}
}
