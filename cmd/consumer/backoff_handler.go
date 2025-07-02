package consumer

import (
	"base/config"
	"base/internal/domain/constants"
	"base/internal/infrastructures/log"
	"base/internal/interfaces/producer"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

func HandleBackoffRetry(ctx context.Context, msg kafka.Message, producer producer.KafkaProducer, cfg config.KafkaConsumerConfig) error {
	attempt := extractRetryAttempt(msg) + 1
	if attempt > cfg.Retry.MaxRetry {
		// return sendToDLQ(ctx, msg, attempt, producer, cfg.Topics.FinalDLQ)
	}

	// Delay logic based on current topic
	delayMinutes := constants.GetBackoffAttemptDelayMinute()[msg.Topic]
	if delayMinutes > 0 {
		delay := time.Duration(delayMinutes) * time.Minute
		log.Printf("Delaying for %v before retry", delay)
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Publish to next retry topic
	nextTopic := getNextBackoffTopic(msg.Topic, cfg)
	if nextTopic == "" {
		return errors.New("no next topic")
	}

	msg.Headers = upsertHeader(msg.Headers, "x-retry-attempt", strconv.Itoa(attempt))
	return producer.Send(ctx, nextTopic, msg.Key, msg.Value)
}

// extractRetryAttempt returns the attempt count from the message header.
func extractRetryAttempt(msg kafka.Message) int {
	for _, h := range msg.Headers {
		if h.Key == "x-retry-attempt" {
			if val, err := strconv.Atoi(string(h.Value)); err == nil {
				return val
			}
		}
	}
	return 0
}

func upsertHeader(headers []kafka.Header, key, value string) []kafka.Header {
	for i, h := range headers {
		if h.Key == key {
			headers[i].Value = []byte(value)
			return headers
		}
	}
	return append(headers, kafka.Header{Key: key, Value: []byte(value)})
}

// getNextBackoffTopic determines the next topic based on the current one.
func getNextBackoffTopic(current string, cfg config.KafkaConsumerConfig) string {
	switch current {
	case constants.Backoff1stAttempt:
		return constants.Backoff2ndAttempt
	case constants.Backoff2ndAttempt:
		return constants.Backoff3rdAttempt
	case constants.Backoff3rdAttempt:
		return ""
		// return constants.FinalDLQ // or "" if you want to stop here
	default:
		return ""
	}
}

// sendToDLQ sends the message to the final DLQ with an additional header.
func sendToDLQ(ctx context.Context, msg kafka.Message, attempt int, producer producer.KafkaProducer, topic string) error {
	msg.Headers = append(msg.Headers, kafka.Header{
		Key:   "x-final-dlq-attempt",
		Value: []byte(strconv.Itoa(attempt)),
	})
	log.Printf("Sending message to DLQ [%s]", topic)
	return producer.Send(ctx, topic, msg.Key, msg.Value)
}
