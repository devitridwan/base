package kafka

import (
	"base/internal/usecases/interactor"
	"context"

	"github.com/segmentio/kafka-go"
)

type GenericKafkaHandler struct {
	interactor interactor.UserInteractor
	topicLogic func(context.Context, kafka.Message, interactor.UserInteractor) error
}

func (h *GenericKafkaHandler) HandleKafkaMessage(ctx context.Context, msg kafka.Message) error {
	return h.topicLogic(ctx, msg, h.interactor)
}
