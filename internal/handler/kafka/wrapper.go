package kafka

// import (
// 	"context"

// 	"github.com/segmentio/kafka-go"
// )

// type KafkaMessageWrapperHandler struct {
// 	handlerFunc func(context.Context, kafka.Message) error
// }

// func (w *KafkaMessageWrapperHandler) HandleKafkaMessage(ctx context.Context, msg kafka.Message) error {
// 	return w.handlerFunc(ctx, msg)
// }

// // func NewKafkaMessageWrapperHandler(fn func(context.Context, kafka.Message) error) KafkaTopicHandler {
// 	return &KafkaMessageWrapperHandler{
// 		handlerFunc: fn,
// 	}
// }
