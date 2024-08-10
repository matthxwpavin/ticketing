package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

func (c *Client) OrderCreatedPublisher(ctx context.Context) (
	streaming.OrderCreatedPublisher,
	error,
) {
	return createStreamIfNotExist[streaming.OrderCreatedMessage](
		ctx,
		c.conn,
		streaming.OrderCreatedStreamConfig,
	)
}

func (c *Client) OrderCreatedConsumer(
	ctx context.Context,
	errHandler streaming.ConsumeErrorHandler,
	filterSubject string) (
	streaming.OrderCreatedConsumer,
	error,
) {
	return consumer[streaming.OrderCreatedMessage](
		ctx,
		c,
		streaming.OrderCreatedStreamConfig,
		errHandler,
		filterSubject,
	)
}

func (c *Client) OrderCancelledPublisher(ctx context.Context) (
	streaming.OrderCancelledPublisher,
	error,
) {
	return createStreamIfNotExist[streaming.OrderCancelledMessage](
		ctx,
		c.conn,
		streaming.OrderCanceledStreamConfig,
	)
}

func (c *Client) OrderCancelledConsumer(
	ctx context.Context,
	errHandler streaming.ConsumeErrorHandler,
	filterSubject string,
) (
	streaming.OrderCancelledConsumer,
	error,
) {
	return consumer[streaming.OrderCancelledMessage](
		ctx,
		c,
		streaming.OrderCanceledStreamConfig,
		errHandler,
		filterSubject,
	)
}
