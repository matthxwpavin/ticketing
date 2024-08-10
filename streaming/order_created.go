package streaming

import (
	"time"
)

const (
	OrderCreatedSubject1 = "order:created:1"
	OrderCreatedSubject2 = "order:created:2"
)

type OrderCreatedPublisher interface {
	Publisher[OrderCreatedMessage]
}

type OrderCreatedConsumer interface {
	JsonConsumer[OrderCreatedMessage]
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
