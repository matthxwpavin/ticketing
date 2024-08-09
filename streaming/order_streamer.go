package streaming

import (
	"context"
)

// TODO: should go to internal sub-package of the orders.
type OrderStreamer interface {
	OrderCreatedPublisher(context.Context) (OrderCreatedPublisher, error)
	OrderCancelledPublisher(context.Context) (OrderCancelledPublisher, error)
	TicketCreatedConsumer(context.Context, ConsumeErrorHandler) (TicketCreatedConsumer, error)
	TicketUpdatedConsumer(context.Context, ConsumeErrorHandler) (TicketUpdateConsumer, error)
	ExpirationCompletedConsumer(context.Context, ConsumeErrorHandler) (ExpirationCompletedConsumer, error)
}
