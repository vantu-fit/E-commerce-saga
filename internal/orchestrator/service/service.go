package service

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/vantu-fit/saga-pattern/internal/orchestrator/service/entity"
	kafkaClient "github.com/vantu-fit/saga-pattern/pkg/kafka"
)

type Service interface {
	StartTransaction(ctx context.Context, purchase *entity.Purchase) error
	HandleReply(ctx context.Context, msg *kafka.Message) error
}

type service struct {
	producer kafkaClient.Producer
}

func NewService(producer kafkaClient.Producer) Service {
	return &service{
		producer: producer,
	}
}
