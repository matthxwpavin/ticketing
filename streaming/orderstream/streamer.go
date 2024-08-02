package orderstream

import "context"

type Streamer interface {
	OrderCreatedPublisher(context.Context) (CreatedPublisher, error)
	OrderCancelledPublisher(context.Context) (CancelledPublisher, error)
}
