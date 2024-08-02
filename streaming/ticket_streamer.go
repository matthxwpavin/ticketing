package streaming

import "context"

type TicketStreamer interface {
	TicketCreatedPublisher(context.Context) (TicketCreatedPublisher, error)
	TicketUpdatedPublisher(context.Context) (TicketUpdatedPublisher, error)
}
