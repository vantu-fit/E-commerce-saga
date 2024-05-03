package event

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/vantu-fit/saga-pattern/cmd/product/config"
	db "github.com/vantu-fit/saga-pattern/internal/product/db/sqlc"
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
	// 	go h.consumer.ConsumeTopic(ctx, poolSize, common.UpdateProductInventoryGroupID, common.UpdateProductInventoryTopic, h.updateProductInventoryWorker)
	// 	go h.consumer.ConsumeTopic(ctx, poolSize, common.RollbackProductInventoryGroupID, common.RollbackProductInventoryTopic, h.rollbackProductInventoryWorker)
}

func (h *eventHandler) updateProductInventoryWorker(ctx context.Context, r kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch message")
			return
		}

		// TODO: handle message
		_ = m
	}
}

func (h *eventHandler) rollbackProductInventoryWorker(ctx context.Context, r kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch message")
			return
		}
		// TODO: handle message
		_ = m
	}
}
