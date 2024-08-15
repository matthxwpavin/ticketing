package streaming

import (
	"context"

	"github.com/matthxwpavin/ticketing/logging/sugar"
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

type StreamConfig struct {
	Name     string
	Subjects []string
}

type AckFunc func() error

type MessgeHadler func([]byte, AckFunc)

type JsonMessageHandler[T any] func(*T, AckFunc)

type Unsubscriber func()

type ConsumeErrorHandler func(Unsubscriber, error)

type TopicsMessages struct {
	TicketUpdatedMsg *TicketUpdatedMessage
	TicketCreatedMsg *TicketCreatedMessage
}

type AcknowledgeMessage[T any] struct {
	acked   bool
	Message *T
}

func (ack *AcknowledgeMessage[T]) Ack()         { ack.acked = true }
func (ack *AcknowledgeMessage[T]) DidAck() bool { return ack.acked }

type AckTopicsMessages struct {
	TicketUpdatedMsg  *AcknowledgeMessage[TicketUpdatedMessage]
	TicketCreatedMsg  *AcknowledgeMessage[TicketCreatedMessage]
	OrderCreatedMsg   *AcknowledgeMessage[OrderCreatedMessage]
	OrderCancelledMsg *AcknowledgeMessage[OrderCancelledMessage]
}

func DefaultConsumeErrorHandler(ctx context.Context) ConsumeErrorHandler {
	return func(_ Unsubscriber, err error) {
		sugar.FromContext(ctx).Errorw("could not consume", "error", err)
	}
}

var (
	TicketCreatedStreamConfig = &StreamConfig{
		Name:     "ticket:created",
		Subjects: []string{"ticket:created:1"},
	}
	TicketUpdatedStreamConfig = &StreamConfig{
		Name:     "ticket:updated",
		Subjects: []string{"ticket:updated:1"},
	}
	OrderCreatedStreamConfig = &StreamConfig{
		Name:     "order:created",
		Subjects: []string{"order:created:1", "order:created:2", "order:created:3"},
	}
	OrderCanceledStreamConfig = &StreamConfig{
		Name:     "order:canceled",
		Subjects: []string{"order:canceled:1", "order:canceled:2", "order:canceled:3"},
	}
	ExpirationCompletedStreamConfig = &StreamConfig{
		Name:     "expiration:completed",
		Subjects: []string{"expiration:completed:1"},
	}
	PaymentCreatedStreamConfig = &StreamConfig{
		Name:     "payment:created",
		Subjects: []string{"payment:created:1"},
	}
)
