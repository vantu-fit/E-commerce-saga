package event

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/vantu-fit/saga-pattern/cmd/mail/config"
	"github.com/vantu-fit/saga-pattern/internal/mail/service"
	"github.com/vantu-fit/saga-pattern/internal/mail/service/command"
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
	producer kafkaClient.Producer
	service  service.Service
}

func NewEvantHandler(
	cfg *config.Config,
	consumer kafkaClient.ConsumerGroup,
	producer kafkaClient.Producer,
	service service.Service,
) EventHanlder {
	return &eventHanlder{
		config:   cfg,
		consumer: consumer,
		producer: producer,
		service:  service,
	}
}

func (h *eventHanlder) Run(ctx context.Context) {
	go h.consumer.ConsumeTopic(ctx, PoolSize, event.SendRegisterEmailGroupID, event.SendRegisterEmailTopic, h.sendRegisterEmailWorker)
}

func (h *eventHanlder) sendRegisterEmailWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Error().Msgf("Mail.SendRegisterEmailWorler: FetchMessage, err: %s", err)
			continue
		}

		var pbEmailRequest pb.SendRegisterEmailRequest
		if err := json.Unmarshal(m.Value, &pbEmailRequest); err != nil {
			log.Error().Msgf("Mail.SendRegisterEmailWorler: Unmarshal, err: %s", err)
			continue
		}

		log.Info().Msgf("Mail.SendRegisterEmailWorler: %v, message at topic/partition/offset %s/%v/%v: %s = %s\n/", workerID, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		reply := pb.SendRegisterEmailResponse{}

		if err = retry.Do(func() error {
			err := h.service.SendRegisterEmailHandler.Handle(ctx , command.Email{
				From: pbEmailRequest.From,
				To: pbEmailRequest.To,
				Subject: pbEmailRequest.Subject,
				FromName: pbEmailRequest.FromName,
				Data: pbEmailRequest.To,
			})
			if err != nil {
				reply.Message = "Send email failed, err: " + err.Error()
				reply.Success = false
			} else {
				reply.Message = "Send email success"
				reply.Success = true
			}

			

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
						Value: []byte(event.SendRegisterEmailTopic),
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
