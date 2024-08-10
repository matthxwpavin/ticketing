package streaming

import (
	"context"
)

// TODO: should go to internal sub-package of the orders.
type OrderStreamer interface {
	OrderCreatedPublisher(context.Context) (OrderCreatedPublisher, error)
	OrderCancelledPublisher(context.Context) (OrderCancelledPublisher, error)
	TicketCreatedConsumer(context.Context, ConsumeErrorHandler, string) (TicketCreatedConsumer, error)
	TicketUpdatedConsumer(context.Context, ConsumeErrorHandler, string) (TicketUpdateConsumer, error)
	ExpirationCompletedConsumer(context.Context, ConsumeErrorHandler, string) (ExpirationCompletedConsumer, error)
}
