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
	OrderId        string    `json:"orderId"`
	OrderStatus    string    `json:"orderStatus"`
	OrderVersion   int32     `json:"version"`
	OrderExpiresAt time.Time `json:"expiresAt"`
	Ticket         struct {
		Id    string `json:"id"`
		Price int32  `json:"price"`
	} `json:"ticket"`
	OrderUserId string `json:"user_id"`
}
