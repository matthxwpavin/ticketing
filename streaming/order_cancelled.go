package streaming

import (
	"context"
)

type OrderCancelledPublisher interface {
	Publish(context.Context, *OrderCancelledMessage) error
}

type OrderCancelledConsumer interface {
	Consume(context.Context, JsonMessageHandler[OrderCancelledMessage]) (Unsubscriber, error)
}

type OrderCancelledMessage struct {
	OrderId string `json:"orderId"`
	Ticket  struct {
		Id string `json:"id"`
	} `json:"ticket"`
}
