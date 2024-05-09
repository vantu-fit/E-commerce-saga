package kafka

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

// Worker kafka consumer worker fetch and process messages from reader
type Worker func(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int)

type ConsumerGroup interface {
	ConsumeTopic(ctx context.Context, poolSize int, groupID string, topic string, worker Worker)
	GetNewKafkaReader(kafkaURL []string, topic, groupID string) *kafka.Reader
	GetNewKafkaWriter() *kafka.Writer
}

type consumerGroup struct {
	Brokers []string
}

func NewConsumerGroup(brokers []string) ConsumerGroup {
	return &consumerGroup{Brokers: brokers}
}

func (c *consumerGroup) GetNewKafkaReader(kafkaURL []string, topic, groupID string) *kafka.Reader {
	return NewKafkaReader(kafkaURL, topic, groupID)
}

func (c *consumerGroup) GetNewKafkaWriter() *kafka.Writer {
	return NewKafkaWriter(c.Brokers)
}

func (c *consumerGroup) ConsumeTopic(ctx context.Context, poolSize int, groupID string, topic string, worker Worker) {
	r := c.GetNewKafkaReader(c.Brokers, topic, groupID)

	defer func() {
		if err := r.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close kafka reader")
		}
	}()

	log.Info().Msgf("Comsumer: %s - Start consuming topic: %s" ,groupID, topic)

	wg := &sync.WaitGroup{}
	for i := 0; i <= poolSize; i++ {
		wg.Add(1)
		go worker(ctx, r, wg, i)
	}
	wg.Wait()
}
