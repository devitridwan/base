package consumer

import (
	"base/internal/domain/repository"
	"base/internal/interfaces/producer"
	"base/internal/interfaces/txmanager"
	"base/internal/usecases/interactor"
)

type module struct {
	userRepository repository.UserRepository
	txManager      txmanager.TxManager
	producer       producer.KafkaProducer
}

type Opts struct {
	UserRepository repository.UserRepository
	TxManager      txmanager.TxManager
	Producer       producer.KafkaProducer
}

func New(o *Opts) interactor.ConsumerInteractor {
	return &module{
		userRepository: o.UserRepository,
		txManager:      o.TxManager,
		producer:       o.Producer,
	}

}
