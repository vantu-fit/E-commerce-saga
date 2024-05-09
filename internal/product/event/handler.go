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
	"github.com/vantu-fit/saga-pattern/cmd/product/config"
	"github.com/vantu-fit/saga-pattern/internal/product/service"
	"github.com/vantu-fit/saga-pattern/internal/product/service/command"
	"github.com/vantu-fit/saga-pattern/internal/product/service/entity"
	"github.com/vantu-fit/saga-pattern/pb"
	"github.com/vantu-fit/saga-pattern/pkg/event"
	kafkaClient "github.com/vantu-fit/saga-pattern/pkg/kafka"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	log.Info().Msg("Start Product event handler")
	go h.consumer.ConsumeTopic(ctx, poolSize, event.UpdateProductInventoryGroupID, event.UpdateProductInventoryTopic, h.updateProductInventoryWorker)
	go h.consumer.ConsumeTopic(ctx, poolSize, event.RollbackProductInventoryGroupID, event.RollbackProductInventoryTopic, h.rollbackProductInventoryWorker)
}

func (h *eventHandler) updateProductInventoryWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Error().Msgf("Product.UpdateProductInventory: FetchMessage, err: %s", err)
			continue
		}

		var purchaseRequest pb.CreatePurchaseRequest
		if err = json.Unmarshal(m.Value, &purchaseRequest); err != nil {
			log.Error().Msgf("Product.UpdateProductInventory: UnmarshalProto , err %s", err)
			continue
		}

		log.Info().Msgf("Product.UpdateProductInventory: %v, message at topic/partition/offset %s/%v/%v: %s = %s\n/", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		reply := pb.CreatePurchaseResponse{
			PurchaseId: purchaseRequest.PurchaseId,
		}

		if err = retry.Do(func() error {
			productItems := make([]entity.ProductItem, len(purchaseRequest.Purchase.Order.OrderItems))
			for i, item := range purchaseRequest.Purchase.Order.OrderItems {
				productItems[i] = entity.ProductItem{
					ID:       uuid.MustParse(item.ProductId),
					Quantity: int64(item.Quantity),
				}
			}

			err = h.service.Command.UpdateProductInventory.Handle(ctx, command.UpdateProductInventory{
				PurchaseID:   uuid.MustParse(purchaseRequest.PurchaseId),
				ProductItems: &productItems,
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
						Value: []byte(event.UpdateProductInventoryHandler),
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
			log.Error().Msgf("Product.UpdateProductInventory: CommitMessage, err: %s", err)
		}

	}
}

func (h *eventHandler) rollbackProductInventoryWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Error().Msgf("Product.RollbackProductInventory: FetchMessage, err: %s", err)
			continue
		}

		var purchaseRequest pb.CreatePurchaseRequest
		if err = json.Unmarshal(m.Value, &purchaseRequest); err != nil {
			log.Error().Msgf("Product.RollbackProductInventory: UnmarshalProto , err %s", err)
			continue
		}

		log.Info().Msgf("Product.RollbackProductInventory: %v, message at topic/partition/offset %s/%v/%v: %s = %s\n/", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		reply := pb.CreatePurchaseResponse{
			PurchaseId: purchaseRequest.PurchaseId,
		}

		if err = retry.Do(func() error {
			productItems := make([]entity.ProductItem, len(purchaseRequest.Purchase.Order.OrderItems))
			for i, item := range purchaseRequest.Purchase.Order.OrderItems {
				productItems[i] = entity.ProductItem{
					ID:       uuid.MustParse(item.ProductId),
					Quantity: int64(item.Quantity),
				}
			}

			err = h.service.Command.RollbackProductInventory.Handle(ctx, command.RollbackProductInventory{
				PurchaseID:   uuid.MustParse(purchaseRequest.PurchaseId),
				ProductItems: &productItems,
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
						Value: []byte(event.RollbackProductInventoryHandler),
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
			log.Error().Msgf("Product.RollbackProductInventory: CommitMessage, err: %s", err)
		}

	}
}
