package consumer

import (
	"base/config"
	kafkaHandler "base/internal/handler/kafka"
	"base/internal/infrastructures/database"
	kafkaInfra "base/internal/infrastructures/kafka"
	"base/internal/infrastructures/log"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

var (
	serveHTTPCmd = &cobra.Command{
		Use:   "run-consumer",
		Short: "Base Service Consumer Kafka",
		Long:  "Serve Base Service Consumer through Kafka",
		RunE:  run,
	}
)

func ServeConsumerCmd() *cobra.Command {
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

	producer, err := kafkaInfra.NewProducer(cfg.Kafka.Producer)
	defer producer.Close()

	appContainer := newContainer(&Opts{
		Cfg:          cfg,
		MasterDataDB: dataStore,
		Producer:     producer,
	})
	handlers := kafkaHandler.NewTopicHandlerMap(cfg.Kafka.Topics, appContainer.ConsumerInteractor)

	dispatcher := kafkaHandler.NewKafkaHandler(kafkaHandler.Opts{
		Handler:            handlers,
		ConsumerInteractor: appContainer.ConsumerInteractor,
	})

	consumer := kafkaInfra.NewConsumer(&kafkaInfra.Opts{
		Cfg: cfg.Kafka,
		Topics: []string{
			cfg.Kafka.Topics.Topic1,
			cfg.Kafka.Topics.Topic2,
		},
		Handler:  dispatcher.Handle,
		Producer: producer,
	})

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Kafka consumer started")
	consumer.Start(ctx)

	<-ctx.Done()

	// Context is canceled → shutdown logic
	log.Println("Shutting down Kafka consumer...")
	consumer.Stop()
	log.Println("Application exited cleanly")

	return nil
}

func dispatchByTopic(
	handlers map[string]func(context.Context, kafka.Message) error,
) func(context.Context, kafka.Message) error {
	return func(ctx context.Context, msg kafka.Message) error {
		handler, ok := handlers[msg.Topic]
		if !ok {
			log.Printf("⚠️ No handler registered for topic: %s", msg.Topic)
			return nil
		}
		return handler(ctx, msg)
	}
}
