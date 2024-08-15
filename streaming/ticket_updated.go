package streaming

type TicketUpdatedPublisher interface {
	Publisher[TicketUpdatedMessage]
}

type TicketUpdateConsumer interface {
	JsonConsumer[TicketUpdatedMessage]
}

type TicketUpdatedMessage struct {
	TicketID      string `json:"ticketID"`
	TicketTitle   string `json:"ticketTitle"`
	TicketPrice   int32  `json:"ticketPrice"`
	TicketVersion int32  `json:"ticketVersion"`
	OrderId       string `json:"orderId"`
}
