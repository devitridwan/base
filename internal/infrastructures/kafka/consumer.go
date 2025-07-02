package kafka

import (
	"base/config"
	"base/internal/domain/constants"
	"base/internal/interfaces/producer"
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader   *kafka.Reader
	handler  func(context.Context, kafka.Message) error
	cfg      config.KafkaConfig
	producer producer.KafkaProducer
}

type Opts struct {
	Cfg      config.KafkaConfig
	Topics   []string
	Handler  func(context.Context, kafka.Message) error
	Producer producer.KafkaProducer
}

func NewConsumer(o *Opts) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     o.Cfg.Consumer.Brokers,
		GroupID:     o.Cfg.Consumer.ConsumerGroup,
		GroupTopics: o.Topics,
	})

	return &Consumer{
		reader:   r,
		handler:  o.Handler,
		cfg:      o.Cfg,
		producer: o.Producer,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Printf("Kafka consumer started with topics: %v", c.reader.Config().GroupTopics)

	const workerCount = 10                 // 🔧 or make this configurable
	jobs := make(chan kafka.Message, 1000) // buffered to allow burst input

	// Start fixed number of workers
	for i := 0; i < workerCount; i++ {
		go c.worker(ctx, jobs)
	}

	// Read messages and feed into job channel
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			return
		}

		select {
		case jobs <- msg:
		case <-ctx.Done():
			log.Println("Kafka consumer shutting down")
			return
		}
	}
}

func (c *Consumer) Stop() {
	log.Println("Closing Kafka reader...")
	_ = c.reader.Close()
}

func (c *Consumer) worker(ctx context.Context, jobs <-chan kafka.Message) {
	for {
		select {
		case msg := <-jobs:
			c.processWithRetry(ctx, msg)
		case <-ctx.Done():
			log.Println("Worker shutdown")
			return
		}
	}
}

func (c *Consumer) processWithRetry(parentCtx context.Context, msg kafka.Message) {
	success := false

	for attempt := 0; attempt <= c.cfg.Consumer.MaxAttempt; attempt++ {
		handlerCtx, cancel := context.WithTimeout(parentCtx, c.cfg.Consumer.Retry.HandlerTimeout)

		err := c.handler(handlerCtx, msg)
		cancel()

		if err == nil {
			success = true
			break
		}

		log.Printf("Handler error (attempt %d/%d): %v", attempt+1, c.cfg.Consumer.MaxAttempt+1, err)

		if attempt < c.cfg.Consumer.MaxAttempt {
			backoff := c.retryDelay(attempt)
			log.Printf("Retrying in %s...", backoff)

			select {
			case <-time.After(backoff):
			case <-parentCtx.Done():
				log.Println("Retry aborted due to shutdown")
				return
			}
		}
	}

	if !success {
		log.Printf("Failed after max retries: %s", string(msg.Value))
		c.sendToBackoffTopic(parentCtx, msg, c.cfg.Consumer.MaxAttempt)
		// Optionally: forward to DLQ here
	}
}

func (c *Consumer) retryDelay(attempt int) time.Duration {
	base := c.cfg.Consumer.Retry.RetryInitialDelay
	if c.cfg.Consumer.Retry.DisableBackoff {
		jitter := time.Duration(rand.Int63n(int64(c.cfg.Consumer.Retry.MaxJitter)))
		return base + jitter
	}

	if d, ok := constants.RetrySleepMap[attempt]; ok {
		return d
	}

	return constants.DefaultRetrySleep

}

func (c *Consumer) sendToBackoffTopic(ctx context.Context, msg kafka.Message, attempt int) {
	topic, _ := constants.GetBackoffAttemptTopicName()[attempt]

	// Clone and enrich the original message
	failedMsg := kafka.Message{
		Topic: topic,
		Key:   msg.Key,
		Value: msg.Value,
		Headers: append(msg.Headers, []kafka.Header{
			{Key: "x-retry-attempt", Value: []byte(strconv.Itoa(attempt))},
			{Key: "x-failed-at", Value: []byte(time.Now().Format(time.RFC3339))},
			{Key: "x-original-topic", Value: []byte(msg.Topic)},
		}...),
	}

	err := c.producer.Send(ctx, topic, failedMsg.Key, failedMsg.Value)
	if err != nil {
		log.Printf("Failed to send to backoff topic [%s]: %v", topic, err)
	} else {
		log.Printf("Sent message to backoff topic [%s] after %d attempts", topic, attempt)
	}
}
