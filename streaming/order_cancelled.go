package streaming

type OrderCancelledPublisher interface {
	Publisher[OrderCancelledMessage]
}

type OrderCancelledConsumer interface {
	JsonConsumer[OrderCancelledMessage]
}

type OrderCancelledMessage struct {
	OrderId      string `json:"orderId"`
	OrderVersion int32  `json:"orderVersion"`
	Ticket       struct {
		Id string `json:"id"`
	} `json:"ticket"`
}
