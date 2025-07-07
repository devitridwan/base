package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type TopicsMap struct {
	Handler map[string]func(ctx context.Context, msg kafka.Message) error
	Topics  []string
}

func NewTopicHandlerMap(handlerReq map[string]func(ctx context.Context, msg kafka.Message) error) TopicsMap {
	topics := []string{}
	for topic, _ := range handlerReq {
		topics = append(topics, topic)
	}
	return TopicsMap{
		Handler: handlerReq,
		Topics:  topics,
	}
}
