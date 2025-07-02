package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaTopicHandler interface {
	HandleKafkaMessage(ctx context.Context, msg kafka.Message) error
}

func NewKafkaHandler(handlers map[string]KafkaTopicHandler) func(ctx context.Context, msg kafka.Message) error {
	return func(ctx context.Context, msg kafka.Message) error {
		handler, ok := handlers[msg.Topic]
		if !ok {
			log.Printf("Unhandled topic: %s", msg.Topic)
			return nil
		}
		return handler.HandleKafkaMessage(ctx, msg)
	}
}
