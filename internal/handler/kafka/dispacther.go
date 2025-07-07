package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaTopicHandler interface {
	HandleKafkaMessage(ctx context.Context, msg kafka.Message) error
}
