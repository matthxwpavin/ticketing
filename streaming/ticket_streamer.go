package streaming

import (
	"context"
)

// TODO: go to internal sub-package of the tickets.
type TicketStreamer interface {
	TicketCreatedPublisher(context.Context) (TicketCreatedPublisher, error)
	TicketUpdatedPublisher(context.Context) (TicketUpdatedPublisher, error)
	OrderCreatedConsumer(context.Context, ConsumeErrorHandler, ...string) (OrderCreatedConsumer, error)
	OrderCancelledConsumer(context.Context, ConsumeErrorHandler, ...string) (OrderCancelledConsumer, error)
	ExpirationCompletedConsumer(context.Context, ConsumeErrorHandler, ...string) (ExpirationCompletedConsumer, error)
}
