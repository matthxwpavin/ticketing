package ticketstream

import "context"

type CreatedPublisher interface {
	Publish(context.Context, *CreatedMessage) error
}

type CreatedMessage struct {
	TicketID      string  `json:"ticketID"`
	TicketTitle   string  `json:"ticketTitle"`
	TicketPrice   float64 `json:"ticketPrice"`
	UserID        string  `json:"userID"`
	TicketVersion int     `json:"ticketVersion"`
}
