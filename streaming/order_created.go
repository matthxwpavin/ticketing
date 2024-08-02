package streaming

import (
	"context"
	"time"
)

type OrderCreatedPublisher interface {
	Publish(context.Context, *OrderCreatedMessage) error
}

type OrderCreatedConsumer interface {
	Consume(context.Context, JsonMessageHandler[OrderCreatedMessage]) (Unsubscriber, error)
}

type OrderCreatedMessage struct {
	OrderId   string    `json:"orderId"`
	Status    string    `json:"orderStatus"`
	ExpiresAt time.Time `json:"expiresAt"`
	Ticket    struct {
		Id    string  `json:"id"`
		Price float64 `json:"price"`
	} `json:"ticket"`
}
