package kafka

import (
	"base/internal/usecases/interactor"
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaTopicHandler interface {
	HandleKafkaMessage(ctx context.Context, msg kafka.Message) error
}

type Opts struct {
	Handler            map[string]KafkaTopicHandler
	ConsumerInteractor interactor.ConsumerInteractor
}

type KafkaHandler struct {
	handler            func(ctx context.Context, msg kafka.Message) error
	consumerInteractor interactor.ConsumerInteractor
}

func NewKafkaHandler(o Opts) *KafkaHandler {
	return &KafkaHandler{
		handler: func(ctx context.Context, msg kafka.Message) error {
			handler, ok := o.Handler[msg.Topic]
			if !ok {
				log.Printf("[WARN] Unhandled topic: %s", msg.Topic)
				return nil
			}
			return handler.HandleKafkaMessage(ctx, msg)
		},
		consumerInteractor: o.ConsumerInteractor,
	}
}

func (h *KafkaHandler) Handle(ctx context.Context, msg kafka.Message) error {
	return h.handler(ctx, msg)
}
