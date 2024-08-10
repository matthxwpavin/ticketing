package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

func (c *Client) orderCreatedSubject() *subject[streaming.OrderCreatedMessage] {
	return &subject[streaming.OrderCreatedMessage]{
		names:           []string{streaming.OrderCreatedSubject1, streaming.OrderCreatedSubject2},
		streamName:      "order:created",
		consumerName:    c.ConsumerName,
		consumerSubject: c.ConsumerSubject,
	}
}

func (c *Client) orderCanceledSubject() *subject[streaming.OrderCancelledMessage] {
	return &subject[streaming.OrderCancelledMessage]{
		names:           []string{"order:canceled"},
		streamName:      "order:canceled",
		consumerName:    c.ConsumerName,
		consumerSubject: c.ConsumerSubject,
	}
}

func (c *Client) OrderCreatedPublisher(ctx context.Context) (
	streaming.OrderCreatedPublisher,
	error,
) {
	return c.orderCreatedSubject().publisher(ctx, c.conn)
}

func (c *Client) OrderCreatedConsumer(ctx context.Context, errHandler streaming.ConsumeErrorHandler) (
	streaming.OrderCreatedConsumer,
	error,
) {
	return c.orderCreatedSubject().jsonConsumer(ctx, c.conn, errHandler)
}

func (c *Client) OrderCancelledPublisher(ctx context.Context) (
	streaming.OrderCancelledPublisher,
	error,
) {
	return c.orderCanceledSubject().publisher(ctx, c.conn)
}

func (c *Client) OrderCancelledConsumer(ctx context.Context, errHandler streaming.ConsumeErrorHandler) (
	streaming.OrderCancelledConsumer,
	error,
) {
	return c.orderCanceledSubject().jsonConsumer(ctx, c.conn, errHandler)
}
