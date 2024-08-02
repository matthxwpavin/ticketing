package streaming

import "context"

type TicketUpdatedPublisher interface {
	Publish(context.Context, *TicketUpdatedMessage) error
}

type TicketUpdatedMessage struct {
	TicketID      string  `json:"ticketID"`
	TicketTitle   string  `json:"ticketTitle"`
	TicketPrice   float64 `json:"ticketPrice"`
	UserID        string  `json:"userID"`
	TicketVersion int     `json:"ticketVersion"`
}
