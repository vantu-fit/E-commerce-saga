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
	"github.com/vantu-fit/saga-pattern/cmd/payment/config"
	"github.com/vantu-fit/saga-pattern/internal/payment/service"
	"github.com/vantu-fit/saga-pattern/internal/payment/service/command"
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
	go h.consumer.ConsumeTopic(ctx, poolSize, event.CreatePaymentGroupID, event.CreatePaymentTopic, h.createPaymentWorker)
	go h.consumer.ConsumeTopic(ctx, poolSize, event.RollbackPaymentGroupID, event.RollbackPaymentTopic, h.rollbackCreatePaymentWorker)
}

func (h *eventHandler) createPaymentWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Error().Msgf("Payment.CreatePaymentWorker: FetchMessage, err: %s", err)
			continue
		}

		var purchaseRequest pb.CreatePurchaseRequest
		if err = json.Unmarshal(m.Value, &purchaseRequest); err != nil {
			log.Error().Msgf("Payment.CreatePaymentWorker: UnmarshalProto , err %s", err)
			continue
		}

		log.Info().Msgf("Payment.CreatePaymentWorker: %v, message at topic/partition/offset %s/%v/%v: %s = %s\n/", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		reply := pb.CreatePurchaseResponse{
			PurchaseId: purchaseRequest.PurchaseId,
		}

		if err = retry.Do(func() error {
			err = h.service.Command.CreatePayment.Handle(ctx, command.CreatePayment{
				ID:         uuid.MustParse(purchaseRequest.PurchaseId),
				CustomerID: uuid.MustParse(purchaseRequest.Purchase.Order.CustomerId),
				Currency:   purchaseRequest.Purchase.Payment.CurrencyCode,
				Amount:     int64(purchaseRequest.Purchase.Payment.Amount),
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
						Value: []byte(event.CreatePaymentHandler),
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
			log.Error().Msgf("Payment.CreatePaymentWorker: CommitMessage, err: %s", err)
		}

	}
}

func (h *eventHandler) rollbackCreatePaymentWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Error().Msgf("Payment.RollbackCreatePaymentWorker: FetchMessage, err: %s", err)
			continue
		}

		var purchaseRequest pb.CreatePurchaseRequest
		if err = json.Unmarshal(m.Value, &purchaseRequest); err != nil {
			log.Error().Msgf("Payment.RollbackCreatePaymentWorker: UnmarshalProto , err %s", err)
			continue
		}

		log.Info().Msgf("Payment.RollbackCreatePaymentWorker: %v, message at topic/partition/offset %s/%v/%v: %s = %s\n/", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		reply := pb.CreatePurchaseResponse{
			PurchaseId: purchaseRequest.PurchaseId,
		}

		if err = retry.Do(func() error {
			err = h.service.Command.DeletePayment.Handle(ctx, command.DeletePayment{
				ID: uuid.MustParse(purchaseRequest.PurchaseId),
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
						Value: []byte(event.RollbackPaymentHandler),
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
			log.Error().Msgf("Payment.RollbackCreatePaymentWorker: CommitMessage, err: %s", err)
		}

	}
}
