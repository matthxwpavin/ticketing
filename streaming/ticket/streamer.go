package ticket

import "context"

type Streamer interface {
	TicketCreatedPublisher(context.Context) (CreatedPublisher, error)
	TicketUpdatedPublisher(context.Context) (UpdatedPublisher, error)
}
