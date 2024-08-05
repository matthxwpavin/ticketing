package streaming

import (
	"context"
)

// Should go to internal sub-package of the tickets.
type TicketStreamer interface {
	TicketCreatedPublisher(context.Context) (TicketCreatedPublisher, error)
	TicketUpdatedPublisher(context.Context) (TicketUpdatedPublisher, error)
	OrderCreatedConsumer(context.Context, ConsumeErrorHandler) (OrderCreatedConsumer, error)
	OrderCancelledConsumer(context.Context, ConsumeErrorHandler) (OrderCancelledConsumer, error)
}
