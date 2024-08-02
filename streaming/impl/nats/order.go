package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

var (
	orderCreated = &subject[streaming.OrderCreatedMessage]{
		name:       "order:created",
		streamName: "order:created",
	}
	orderCancelled = &subject[streaming.OrderCancelledMessage]{
		name:       "order:cancelled",
		streamName: "order:cancelled",
	}
)

func (c *Client) OrderCreatedPublisher(ctx context.Context) (
	streaming.OrderCreatedPublisher,
	error,
) {
	return orderCreated.publisher(ctx, c.conn)
}

func (c *Client) OrderCreatedConsumer(ctx context.Context, errHandler streaming.ConsumeErrorHandler) (
	streaming.OrderCreatedConsumer,
	error,
) {
	return orderCreated.jsonConsumer(ctx, c.conn, errHandler)
}

func (c *Client) OrderCancelledPublisher(ctx context.Context) (
	streaming.OrderCancelledPublisher,
	error,
) {
	return orderCancelled.publisher(ctx, c.conn)
}

func (c *Client) OrderCancelledConsumer(ctx context.Context, errHandler streaming.ConsumeErrorHandler) (
	streaming.OrderCancelledConsumer,
	error,
) {
	return orderCancelled.jsonConsumer(ctx, c.conn, errHandler)
}
