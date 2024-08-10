package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

func (c *Client) expirationCompletedSubject() *subject[streaming.ExpirationCompletedMessage] {
	return &subject[streaming.ExpirationCompletedMessage]{
		name:         "expiration:completed",
		streamName:   "expiration:completed",
		consumerName: c.ConsumerName,
	}
}

func (c *Client) ExpirationCompletedConsumer(ctx context.Context, errHandler streaming.ConsumeErrorHandler) (
	streaming.ExpirationCompletedConsumer,
	error,
) {
	return c.expirationCompletedSubject().jsonConsumer(ctx, c.conn, errHandler)
}

func (c *Client) ExpirationCompletedPublisher(ctx context.Context) (
	streaming.ExpirationCompletedPublisher,
	error,
) {
	return c.expirationCompletedSubject().publisher(ctx, c.conn)
}
