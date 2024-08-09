package streaming

import (
	"time"
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
