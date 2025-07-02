package kafka

import (
	"base/config"
	"base/internal/interfaces/producer"
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewProducer(cfg config.KafkaProducerConfig) (producer.KafkaProducer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("no Kafka brokers configured")
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireAll,
	}

	if cfg.MaxMessageSize > 0 {
		writer.MaxAttempts = 3
		writer.WriteTimeout = 10 * time.Second
	}

	return &kafkaProducer{
		writer: writer,
	}, nil
}

func (p *kafkaProducer) Send(ctx context.Context, topic string, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}

func (p *kafkaProducer) Close() error {
	return p.writer.Close()
}
