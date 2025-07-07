package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type TopicHandler struct {
	handler func(ctx context.Context, msg kafka.Message) error
}
