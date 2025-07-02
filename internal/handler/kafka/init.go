package kafka

import (
	"base/config"
	"base/internal/usecases/interactor"
)

func NewTopicHandlerMap(cfg config.KafkaTopics, interactor interactor.ConsumerInteractor) map[string]KafkaTopicHandler {
	return map[string]KafkaTopicHandler{
		cfg.Topic1: &UnifiedKafkaHandler{
			Interactor: interactor,
			Cfg:        cfg,
		},
		cfg.Topic2: &UnifiedKafkaHandler{
			Interactor: interactor,
			Cfg:        cfg,
		},
		// Add more topics here
	}
}
