package event

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/vantu-fit/saga-pattern/cmd/purchase/config"
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
	ProduceCreatePurchaseEvent(ctx context.Context, purchase *pb.CreatePurchaseRequest) error
	Run(ctx context.Context)
}

type eventHanlder struct {
	config   *config.Config
	producer kafkaClient.Producer
	consumer kafkaClient.ConsumerGroup
}

func NewEvantHandler(
	cfg *config.Config, 
	producer kafkaClient.Producer,
	consumer kafkaClient.ConsumerGroup,
) EventHanlder {
	return &eventHanlder{
		config:   cfg,
		producer: producer,
		consumer: consumer,
	}
}

func (h *eventHanlder) ProduceCreatePurchaseEvent(ctx context.Context, purchase *pb.CreatePurchaseRequest) error {
	payload, err := json.Marshal(purchase)
	if err != nil {
		return err
	}

	return h.producer.PublishMessage(ctx, kafka.Message{
		Topic: event.PurchaseTopic,
		Key:   []byte("create_purchase"),
		Value: payload,
	})
}

func (h *eventHanlder) Run(ctx context.Context) {
	go h.consumer.ConsumeTopic(ctx, PoolSize, event.PurchaseGroupID, event.PurchaseResultTopic, h.purchaseresultWorker)
}

func (h *eventHanlder) purchaseresultWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Error().Msgf("Orchestartor.PurchaseResultWorker: FetchMessage, err: %s", err)
			continue
		}

		log.Info().Msgf("Purchase.PurchaseResultWorkder: %v, message at topic/partition/offset %s/%v/%v: %s = %s\n/", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		err = r.CommitMessages(ctx, m)
		if err != nil {
			log.Error().Msgf("Purchase.PurchaseResultWorkder: CommitMessage, err: %s", err)
		}
	}
}
