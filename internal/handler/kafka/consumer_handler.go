package kafka

import (
	"base/config"
	"base/internal/domain/models"
	"base/internal/handler/api/request"
	"base/internal/interfaces/producer"
	"base/internal/usecases/interactor"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaHandler struct {
	consumerInteractor interactor.ConsumerInteractor
	producer           producer.KafkaProducer
	cfg                config.KafkaConsumerConfig
	topics             config.KafkaTopics
}

type KafkaHandlerOpts struct {
	ConsumerInteractor interactor.ConsumerInteractor
	Producer           producer.KafkaProducer
	Cfg                config.KafkaConsumerConfig
	Topics             config.KafkaTopics
}

func NewKafkaHandler(o KafkaHandlerOpts) *KafkaHandler {
	return &KafkaHandler{
		consumerInteractor: o.ConsumerInteractor,
		producer:           o.Producer,
		cfg:                o.Cfg,
		topics:             o.Topics,
	}
}

func (h *KafkaHandler) Topic1(ctx context.Context, msg kafka.Message) (err error) {
	var req request.GetByIdReq
	fmt.Println(string(msg.Value))
	if err = json.Unmarshal(msg.Value, &req); err != nil {
		log.Printf("[ERROR] Invalid message JSON: %v", err)
		return err
	}
	_, err = h.consumerInteractor.Topic1(ctx, &req)
	return err
}

func (h *KafkaHandler) Topic2(ctx context.Context, msg kafka.Message) (err error) {
	var req request.GetByIdReq
	if err = json.Unmarshal(msg.Value, &req); err != nil {
		log.Printf("[ERROR] Invalid message JSON: %v", err)
		return err
	}
	_, err = h.consumerInteractor.Topic2(ctx, &req)
	return err
}

func (h *KafkaHandler) Backoff(ctx context.Context, msg kafka.Message) (err error) {
	var req models.BackoffRetry
	if err := json.Unmarshal(msg.Value, &req); err != nil {
		log.Printf("[ERROR] Invalid message JSON: %v", err)
		return err
	}
	origin := msg
	msg.Value = req.Data
	fmt.Println("value >> ", string(msg.Value))
	switch req.SourceTopic {
	case h.topics.Topic1:
		err = h.Topic1(ctx, msg)
	case h.topics.Topic2:
		err = h.Topic2(ctx, msg)
	}
	if err == nil {
		return nil
	}
	return HandleBackoffRetry(ctx, origin, h.producer, h.cfg)
}
