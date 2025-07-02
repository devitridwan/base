package server

import (
	"base/config"

	"base/internal/infrastructures/database"
	"base/internal/infrastructures/kafka"

	_ "github.com/lib/pq"

	"base/internal/handler"
	"base/internal/infrastructures/log"

	"github.com/spf13/cobra"
)

var (
	serveHTTPCmd = &cobra.Command{
		Use:   "serve-http",
		Short: "User Service HTTP",
		Long:  "Serve User Service through HTTP",
		RunE:  run,
	}
)

func ServeHTTPCmd() *cobra.Command {
	serveHTTPCmd.Flags().StringP("config", "c", "", "Config Path, both relative or absolute. i.e: /usr/local/bin/config/files")
	return serveHTTPCmd
}

func run(cmd *cobra.Command, args []string) error {

	configLocation, err := cmd.Flags().GetString("config")
	if err != nil {
		log.WithError(err).Fatalf("Failed to get config: %v", err)
	}

	cfg := &config.MainConfig{}
	config.ReadModuleConfig(cfg, "main", configLocation)

	/* NOT USED YET
	_, err = tracing.GetOtelProvider(cfg.Otel.Kind, constants.PayslipService, cfg.Otel)
	if err != nil {
		log.WithError(err).Fatalf("Failed to initiate jaeger: %v", err)
	}
	*/

	dataStore := database.New(database.DBConfig{
		SlaveDSN:        cfg.Database.SlaveDSN,
		MasterDSN:       cfg.Database.MasterDSN,
		RetryInterval:   cfg.Database.RetryInterval,
		MaxIdleConn:     cfg.Database.MaxIdleConn,
		MaxConn:         cfg.Database.MaxConn,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}, database.DriverPostgres)

	producer, err := kafka.NewProducer(cfg.Kafka.Producer)
	defer producer.Close()

	appContainer := newContainer(&Opts{
		Cfg:          cfg,
		MasterDataDB: dataStore,
		Producer:     producer,
	})

	server := handler.NewHTTP(&handler.Opts{
		Cfg:            appContainer.Cfg,
		UserInteractor: appContainer.UserInteractor,
	})

	server.Run()

	return nil
}
