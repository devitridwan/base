package kafka

import (
	"base/internal/handler/api/request"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func (h *KafkaHandler) Topic1(ctx context.Context, msg kafka.Message) error {
	var req request.GetByIdReq
	if err := json.Unmarshal(msg.Value, &req); err != nil {
		log.Printf("[ERROR] Invalid message JSON: %v", err)
		return err
	}
	_, err := h.consumerInteractor.Topic1(ctx, &req)
	return err
}

func (h *KafkaHandler) Topic2(ctx context.Context, msg kafka.Message) error {
	var req request.GetByIdReq
	if err := json.Unmarshal(msg.Value, &req); err != nil {
		log.Printf("[ERROR] Invalid message JSON: %v", err)
		return err
	}
	_, err := h.consumerInteractor.Topic2(ctx, &req)
	return err
}
