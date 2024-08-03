package streaming

import (
	"context"
)

type TicketCreatedPublisher interface {
	Publish(context.Context, *TicketCreatedMessage) error
}

type TicketCreatedConsumer interface {
	Consume(context.Context, JsonMessageHandler[TicketCreatedMessage]) (Unsubscriber, error)
}

type TicketCreatedMessage struct {
	TicketID      string  `json:"ticketID"`
	TicketTitle   string  `json:"ticketTitle"`
	TicketPrice   float64 `json:"ticketPrice"`
	UserID        string  `json:"userID"`
	TicketVersion int32   `json:"ticketVersion"`
}
