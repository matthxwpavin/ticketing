package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming/orderstream"
)

func (c *Client) OrderCreatedPublisher(ctx context.Context) (
	orderstream.CreatedPublisher,
	error,
) {
	return publisher[orderstream.CreatedMessage](ctx, c.conn, "order:created", "order:created")
}

func (c *Client) OrderCancelledPublisher(ctx context.Context) (
	orderstream.CancelledPublisher,
	error,
) {
	return publisher[orderstream.CancelledMessage](ctx, c.conn, "order:cancelled", "order:cancelled")
}
