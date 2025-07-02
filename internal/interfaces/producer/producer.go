package producer

import "context"

type KafkaProducer interface {
	Send(ctx context.Context, topic string, key, value []byte) error
	Close() error
}
