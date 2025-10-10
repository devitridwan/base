package kafka

import (
	"base/config"
	"base/internal/domain/constants"
	"base/internal/domain/models"
	"base/internal/interfaces/producer"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader   *kafka.Reader
	handlers map[string]func(ctx context.Context, msg kafka.Message) error
	cfg      config.KafkaConfig
	producer producer.KafkaProducer
}

type Opts struct {
	Cfg      config.KafkaConfig
	Topics   []string
	Handlers map[string]func(ctx context.Context, msg kafka.Message) error
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
		handlers: o.Handlers,
		cfg:      o.Cfg,
		producer: o.Producer,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Printf("Kafka consumer started with topics: %v", c.reader.Config().GroupTopics)

	const workerCount = 10
	jobs := make(chan kafka.Message, 1000)

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

func (c *Consumer) processWithRetry(ctx context.Context, msg kafka.Message) {
	var errHandler error

	handler, ok := c.handlers[msg.Topic]
	if !ok {
		log.Printf("[WARN] No handler for topic: %s", msg.Topic)
		return
	}

	success := false

	for attempt := 0; attempt <= c.cfg.Consumer.Retry.MaxRetry; attempt++ {
		ctx, cancel := context.WithTimeout(ctx, c.cfg.Consumer.Retry.HandlerTimeout)

		errHandler := handler(ctx, msg)
		cancel()

		if errHandler == nil {
			success = true
			return
		}

		log.Printf("Handler error (attempt %d/%d): %v", attempt+1, c.cfg.Consumer.Retry.MaxRetry+1, errHandler)

		if attempt < c.cfg.Consumer.Retry.MaxRetry {
			backoff := c.retryDelay(attempt)
			log.Printf("Retrying in %s...", backoff)

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				log.Println("Retry aborted due to shutdown")
				return
			}
		}
	}
	key := 1
	temp, ok := GetHeader(msg.Headers, "x-retry-attempt")
	if ok {
		key = temp
	}

	messValue := models.BackoffRetry{
		SourceTopic: msg.Topic,
		Data:        msg.Value,
		Timestamp:   time.Now(),
		Error:       errHandler.Error(),
	}
	byt, _ := json.Marshal(messValue)

	msg.Value = byt

	if !success {
		log.Printf("Failed after max retries: %s", string(msg.Value))
		c.sendToBackoffTopic(ctx, msg, key)
	}
}

func (c *Consumer) retryDelay(attempt int) time.Duration {
	base := c.cfg.Consumer.Retry.RetryInitialDelay
	if c.cfg.Consumer.Retry.DisableBackoff {
		jitter := time.Duration(rand.Int63n(int64(c.cfg.Consumer.Retry.MaxJitter)))
		return base + jitter
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
		}...),
	}

	err := c.producer.Send(ctx, topic, failedMsg.Key, failedMsg.Value)
	if err != nil {
		log.Printf("Failed to send to backoff topic [%s]: %v", topic, err)
	} else {
		log.Printf("Sent message to backoff topic [%s] after %d attempts", topic, attempt)
	}
}

func GetHeader(headers []kafka.Header, key string) (int, bool) {
	for _, h := range headers {
		if h.Key == key {
			val, _ := strconv.Atoi(string(h.Value))
			return val, true
		}
	}
	return 0, false
}
