package orderstream

import (
	"context"
	"time"
)

type CreatedPublisher interface {
	Publish(context.Context, *CreatedMessage) error
}

type CreatedMessage struct {
	OrderId   string    `json:"orderId"`
	Status    string    `json:"orderStatus"`
	ExpiresAt time.Time `json:"expiresAt"`
	Ticket    struct {
		Id    string  `json:"id"`
		Price float64 `json:"price"`
	} `json:"ticket"`
}
