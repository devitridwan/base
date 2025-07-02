package kafka

import (
	"base/config"
	"base/internal/handler/api/request"
	"base/internal/usecases/interactor"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type UnifiedKafkaHandler struct {
	Interactor interactor.ConsumerInteractor
	Cfg        config.KafkaTopics
}

func (h *UnifiedKafkaHandler) HandleKafkaMessage(ctx context.Context, msg kafka.Message) error {
	switch msg.Topic {
	case h.Cfg.Topic1:
		var req request.GetByIdReq
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			log.Printf("[ERROR] Failed to unmarshal Topic1 message: %v", err)
			return err
		}
		fmt.Println(req)
		_, err := h.Interactor.Topic1(ctx, &req)
		return err
	case h.Cfg.Topic2:
		var req request.GetByIdReq
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			log.Printf("[ERROR] Failed to unmarshal Topic2 message: %v", err)
			return err
		}
		_, err := h.Interactor.Topic2(ctx, &req)
		return err
	default:
		log.Printf("[WARN] No handler for topic: %s", msg.Topic)
		return nil
	}
}
