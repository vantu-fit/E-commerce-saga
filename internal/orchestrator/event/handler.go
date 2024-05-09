package event

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/vantu-fit/saga-pattern/cmd/orchestrator/config"
	"github.com/vantu-fit/saga-pattern/internal/orchestrator/service"
	"github.com/vantu-fit/saga-pattern/internal/orchestrator/service/entity"
	"github.com/vantu-fit/saga-pattern/pb"
	"github.com/vantu-fit/saga-pattern/pkg/event"
	kafkaClient "github.com/vantu-fit/saga-pattern/pkg/kafka"
)

var (
	retryAttempts uint = 10
	retryDelay         = 1 * time.Second
	PoolSize           = 16
)

type EventHanlder interface {
	Run(ctx context.Context)
}

type eventHanlder struct {
	config   *config.Config
	consumer kafkaClient.ConsumerGroup
	service  service.Service
}

func NewEvantHandler(cfg *config.Config, consumer kafkaClient.ConsumerGroup, service service.Service) EventHanlder {
	return &eventHanlder{
		config:   cfg,
		consumer: consumer,
		service:  service,
	}
}

func (h *eventHanlder) Run(ctx context.Context) {
	go h.consumer.ConsumeTopic(ctx, PoolSize, event.PurchaseGroupID, event.PurchaseTopic, h.createPurchaseWorker)
	go h.consumer.ConsumeTopic(ctx, PoolSize, event.ReplyGroupID, event.ReplyTopic, h.replyWorker)
}

func (h *eventHanlder) createPurchaseWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Error().Msgf("Orchestartor.CreatePurchaseWorker: FetchMessage, err: %s", err)
			continue
		}

		var purchaseRequest pb.CreatePurchaseRequest
		if err = json.Unmarshal(m.Value, &purchaseRequest); err != nil {
			log.Error().Msgf("Orchestrator.CreatePurchaseWorker: UnmarshalProto , err %s", err)
			continue
		}

		log.Info().Msgf("Orchestartor.CreatePurchaseWorker: %v, message at topic/partition/offset %s/%v/%v: %s = %s\n/", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		orderItems := make([]entity.OrderItem, len(purchaseRequest.Purchase.Order.OrderItems))
		for i, orderItem := range purchaseRequest.Purchase.Order.OrderItems {
			orderItems[i] = entity.OrderItem{
				ID:       uuid.MustParse(orderItem.ProductId),
				Quantity: orderItem.Quantity,
			}
		}

		purchase := entity.Purchase{
			ID: uuid.MustParse(purchaseRequest.PurchaseId),
			Order: &entity.Order{
				ID:         uuid.MustParse(purchaseRequest.PurchaseId),
				CustomerID: uuid.MustParse(purchaseRequest.Purchase.Order.CustomerId),
				OrderItems: &orderItems,
			},
			Payment: &entity.Payment{
				ID:            uuid.MustParse(purchaseRequest.PurchaseId),
				CurrentcyCode: purchaseRequest.Purchase.Payment.CurrencyCode,
				Amount:        purchaseRequest.Purchase.Payment.Amount,
			},
		}

		if err = retry.Do(func() error {
			err = h.service.StartTransaction(ctx, &purchase)
			if err != nil {
				log.Error().Msgf("Orchestrator.CreatePurchaseWorker: StartTransaction, err: %s", err)
			}
			return err
		},
			retry.Attempts(retryAttempts),
			retry.Delay(retryDelay),
			retry.Context(ctx),
		); err != nil {
			// TODO puplic error message to error topic
		}

		err = r.CommitMessages(ctx, m)
		if err != nil {
			log.Error().Msgf("Orchestrator.CreatePurchaseWorker: CommitMessage, err: %s", err)
		}

	}

}
func (h *eventHanlder) replyWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()


	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Error().Msgf("Orchestrator.ReplyWorker Fetch Message: %s", err)
			continue
		}

		log.Info().Msgf("Orchestartor.ReplyWorker: %v, message at topic/partition/offset %s/%v/%v: %s = %s\n/", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		if err = retry.Do(func() error {
			err = h.service.HandleReply(ctx , &m)
			if err != nil {
				log.Error().Msgf("Orchestrator.ReplyWorker: StartTransaction, err: %s", err)
			}
			return err
		},
		retry.Attempts(retryAttempts),
			retry.Delay(retryDelay),
			retry.Context(ctx),
		); err != nil {
			// TODO puplic error message to error topic
		}

		err = r.CommitMessages(ctx, m)
		if err != nil {
			log.Error().Msgf("Orchestrator.CreatePurchaseWorker: CommitMessage, err: %s", err)
		}
	}
}
