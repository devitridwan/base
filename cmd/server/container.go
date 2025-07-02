package server

import (
	"base/config"
	"base/internal/infrastructures/database"
	"base/internal/interfaces/dao"
	"base/internal/interfaces/producer"
	"base/internal/interfaces/txmanager"
	"base/internal/usecases/interactor"
	user "base/internal/usecases/user"
)

type Container struct {
	Cfg            config.MainConfig
	UserInteractor interactor.UserInteractor
}

type Opts struct {
	Cfg          *config.MainConfig
	MasterDataDB *database.Store
	Producer     producer.KafkaProducer
}

func newContainer(o *Opts) *Container {
	tx := txmanager.NewTxManager(&txmanager.Opts{
		DB: o.MasterDataDB,
	})

	userRepo := dao.NewUserRepo(&dao.OptUser{DB: o.MasterDataDB})

	userInteractor := user.New(&user.Opts{
		UserRepository: userRepo,
		TxManager:      tx,
		Producer:       o.Producer,
	})

	return &Container{
		Cfg:            *o.Cfg,
		UserInteractor: userInteractor,
	}
}
