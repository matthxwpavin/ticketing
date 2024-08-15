package streaming

type TicketCreatedPublisher interface {
	Publisher[TicketCreatedMessage]
}

type TicketCreatedConsumer interface {
	JsonConsumer[TicketCreatedMessage]
}

type TicketCreatedMessage struct {
	TicketID      string `json:"ticketID"`
	TicketTitle   string `json:"ticketTitle"`
	TicketPrice   int32  `json:"ticketPrice"`
	UserID        string `json:"userID"`
	TicketVersion int32  `json:"ticketVersion"`
}
