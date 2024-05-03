package event

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/vantu-fit/saga-pattern/cmd/order/config"
	"github.com/vantu-fit/saga-pattern/internal/common"
	db "github.com/vantu-fit/saga-pattern/internal/order/db/sqlc"
	kafkaClient "github.com/vantu-fit/saga-pattern/pkg/kafka"
)

var (
	retryAttempts uint = 10
	retryDelay         = 1 * time.Second
	poolSize           = 16
)

type EventHandler interface {
	Run(ctx context.Context)
}

type eventHandler struct {
	cfg        *config.Config
	consumer   kafkaClient.ConsumerGroup
	producer   kafkaClient.Producer
	storeCache db.Store
}

func NewEventHandler(
	cfg *config.Config,
	consumer kafkaClient.ConsumerGroup,
	producer kafkaClient.Producer,
	storeCache db.Store,
) EventHandler {
	return &eventHandler{
		cfg:        cfg,
		consumer:   consumer,
		producer:   producer,
		storeCache: storeCache,
	}
}

func (h *eventHandler) Run(ctx context.Context) {
	log.Info().Msg("Start event handler")
	go h.consumer.ConsumeTopic(ctx, poolSize, common.CreateOrderGroupID, common.CreateOrderTopic, h.createOrderWorker)
	go h.consumer.ConsumeTopic(ctx, poolSize, common.RollbackOrderGroupID, common.RollbackOrderTopic, h.rollbackCreateOrderWorker)
}

func (h *eventHandler) createOrderWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int){
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		log.Debug().Msgf(string(m.Value))
		if err != nil {
			log.Error().Err(err).Msg("Order.CreateOrderWorker: FetchMessage")
			return
		}

		// TODO: handle message
		_ = m
	}
}

func (h *eventHandler) rollbackCreateOrderWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		log.Debug().Msgf(string(m.Value))
		if err != nil {
			log.Error().Err(err).Msg("Order.RollbackOrderWorker: FetchMessage")
			return
		}
		// TODO: handle message
		_ = m
	}
}
