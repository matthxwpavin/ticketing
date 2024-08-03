package streaming

import (
	"context"
)

type OrderStreamer interface {
	OrderCreatedPublisher(context.Context) (OrderCreatedPublisher, error)
	OrderCancelledPublisher(context.Context) (OrderCancelledPublisher, error)
	TicketCreatedConsumer(context.Context, ConsumeErrorHandler) (TicketCreatedConsumer, error)
	TicketUpdatedConsumer(context.Context, ConsumeErrorHandler) (TicketUpdateConsumer, error)
}
