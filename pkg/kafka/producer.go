package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishMessage(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type producer struct {
	brokers []string
	w       *kafka.Writer
}

// NewProducer create new kafka producer
func NewProducer(brokers []string) Producer {
	return &producer{brokers: brokers, w: NewKafkaWriter(brokers)}
}

func (p *producer) PublishMessage(ctx context.Context, msgs ...kafka.Message) error {
	if err := p.w.WriteMessages(ctx, msgs...); err != nil {
		return err
	}
	return nil
}

func (p *producer) Close() error {
	return p.w.Close()
}
