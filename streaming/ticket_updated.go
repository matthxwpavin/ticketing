package streaming

import "context"

type TicketUpdatedPublisher interface {
	Publish(context.Context, *TicketUpdatedMessage) error
}

type TicketUpdateConsumer interface {
	Consume(context.Context, JsonMessageHandler[TicketUpdatedMessage]) (Unsubscriber, error)
}

type TicketUpdatedMessage struct {
	TicketID      string  `json:"ticketID"`
	TicketTitle   string  `json:"ticketTitle"`
	TicketPrice   float64 `json:"ticketPrice"`
	TicketVersion int32   `json:"ticketVersion"`
}
