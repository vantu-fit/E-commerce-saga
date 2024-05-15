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
	"github.com/vantu-fit/saga-pattern/cmd/order/config"
	"github.com/vantu-fit/saga-pattern/internal/order/service"
	"github.com/vantu-fit/saga-pattern/internal/order/service/command"
	"github.com/vantu-fit/saga-pattern/pb"
	"github.com/vantu-fit/saga-pattern/pkg/event"
	kafkaClient "github.com/vantu-fit/saga-pattern/pkg/kafka"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	retryAttempts uint = 10
	retryDelay         = 1 * time.Second
	poolSize           = 64
)

type EventHandler interface {
	Run(ctx context.Context)
}

type eventHandler struct {
	cfg      *config.Config
	consumer kafkaClient.ConsumerGroup
	producer kafkaClient.Producer
	service  *service.Service
}

func NewEventHandler(
	cfg *config.Config,
	consumer kafkaClient.ConsumerGroup,
	producer kafkaClient.Producer,
	service *service.Service,
) EventHandler {
	return &eventHandler{
		cfg:      cfg,
		consumer: consumer,
		producer: producer,
		service:  service,
	}
}

func (h *eventHandler) Run(ctx context.Context) {
	log.Info().Msg("Start event handler")
	go h.consumer.ConsumeTopic(ctx, poolSize, event.CreateOrderGroupID, event.CreateOrderTopic, h.createOrderWorker)
	go h.consumer.ConsumeTopic(ctx, poolSize, event.RollbackOrderGroupID, event.RollbackOrderTopic, h.rollbackCreateOrderWorker)
}

func (h *eventHandler) createOrderWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Error().Msgf("Order.CreateOrderWorker: FetchMessage, err: %s", err)
			continue
		}

		var purchaseRequest pb.CreatePurchaseRequest
		if err = json.Unmarshal(m.Value, &purchaseRequest); err != nil {
			log.Error().Msgf("Order.CreateOrderWorker: UnmarshalProto , err %s", err)
			continue
		}

		log.Info().Msgf("Order.CreateOrderWorker: %v, message at topic/partition/offset %s/%v/%v: %s = %s\n/", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		reply := pb.CreatePurchaseResponse{
			PurchaseId: purchaseRequest.PurchaseId,
		}

		purchaseProducts := make([]command.PurchasedProduct, len(purchaseRequest.Purchase.Order.OrderItems))
		for i, purchaseProduct := range purchaseRequest.Purchase.Order.OrderItems {
			purchaseProducts[i] = command.PurchasedProduct{
				ProductID: uuid.MustParse(purchaseProduct.ProductId),
				Quantity:  purchaseProduct.Quantity,
			}
		}

		if err = retry.Do(func() error {
			err = h.service.Command.CreateOrder.Handle(ctx, command.CreateOrder{
				OrderID:    uuid.MustParse(purchaseRequest.PurchaseId),
				CustomerID: uuid.MustParse(purchaseRequest.Purchase.Order.CustomerId),
				Products:   &purchaseProducts,
			})
			if err != nil {
				reply.Success = false
				reply.ErrorMessage = err.Error()
			} else {
				reply.Success = true
			}

			reply.Purchase = purchaseRequest.Purchase
			reply.Timestamp = timestamppb.Now()

			payload, err := json.Marshal(&reply)
			if err != nil {
				return err
			}

			return h.producer.PublishMessage(ctx, kafka.Message{
				Topic: event.ReplyTopic,
				Value: payload,
				Headers: []kafka.Header{
					{
						Key:   event.HandlerHeader,
						Value: []byte(event.CreateOrderHandler),
					},
				},
			})
		},
			retry.Attempts(retryAttempts),
			retry.Delay(retryDelay),
			retry.Context(ctx),
		); err != nil {
			// TODO puplic error message to error topic
		}

		err = r.CommitMessages(ctx, m)
		if err != nil {
			log.Error().Msgf("Order.CreateOrderWorker: CommitMessage, err: %s", err)
		}

	}
}

func (h *eventHandler) rollbackCreateOrderWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Error().Msgf("Order.RollbackOrder: FetchMessage, err: %s", err)
			continue
		}

		var purchaseRequest pb.CreatePurchaseRequest
		if err = json.Unmarshal(m.Value, &purchaseRequest); err != nil {
			log.Error().Msgf("Order.RollbackOrder: UnmarshalProto , err %s", err)
			continue
		}

		log.Info().Msgf("Order.RollbackOrder: %v, message at topic/partition/offset %s/%v/%v: %s = %s\n/", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		reply := pb.CreatePurchaseResponse{
			PurchaseId: purchaseRequest.PurchaseId,
		}

		purchaseProducts := make([]command.PurchasedProduct, len(purchaseRequest.Purchase.Order.OrderItems))
		for i, purchaseProduct := range purchaseRequest.Purchase.Order.OrderItems {
			purchaseProducts[i] = command.PurchasedProduct{
				ProductID: uuid.MustParse(purchaseProduct.ProductId),
				Quantity:  purchaseProduct.Quantity,
			}
		}

		if err = retry.Do(func() error {
			err = h.service.Command.DeleteOrder.Handle(ctx, command.DeleteOrder{
				OrderID: uuid.MustParse(purchaseRequest.PurchaseId),
			})
			if err != nil {
				reply.Success = false
				reply.ErrorMessage = err.Error()
			} else {
				reply.Success = true
			}

			reply.Purchase = purchaseRequest.Purchase
			reply.Timestamp = timestamppb.Now()

			payload, err := json.Marshal(&reply)
			if err != nil {
				return err
			}

			return h.producer.PublishMessage(ctx, kafka.Message{
				Topic: event.ReplyTopic,
				Value: payload,
				Headers: []kafka.Header{
					{
						Key:   event.HandlerHeader,
						Value: []byte(event.RollbackOrderHandler),
					},
				},
			})
		},
			retry.Attempts(retryAttempts),
			retry.Delay(retryDelay),
			retry.Context(ctx),
		); err != nil {
			// TODO puplic error message to error topic
		}

		err = r.CommitMessages(ctx, m)
		if err != nil {
			log.Error().Msgf("Order.RollbackOrder: CommitMessage, err: %s", err)
		}

	}
}
