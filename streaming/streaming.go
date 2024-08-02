package streaming

import (
	"context"
)

type Publisher[T any] interface {
	Publish(ctx context.Context, message *T) error
}

type Consumer interface {
	Consume(context.Context, MessgeHadler) (Unsubscriber, error)
}
type JsonConsumer[T any] interface {
	Consume(context.Context, JsonMessageHandler[T]) (Unsubscriber, error)
}

type AckFunc func() error

type MessgeHadler func([]byte, AckFunc)

type JsonMessageHandler[T any] func(*T, AckFunc)

type Unsubscriber func()

type ConsumeErrorHandler func(Unsubscriber, error)
