package ticketstream

import "context"

type UpdatedPublisher interface {
	Publish(context.Context, *UpdatedMessage) error
}

type UpdatedMessage struct {
	TicketID      string  `json:"ticketID"`
	TicketTitle   string  `json:"ticketTitle"`
	TicketPrice   float64 `json:"ticketPrice"`
	UserID        string  `json:"userID"`
	TicketVersion int     `json:"ticketVersion"`
}
