package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

var expirationCompleted = &subject[streaming.ExpirationCompletedMessage]{
	name:       "expiration:completed",
	streamName: "expiration:completed",
}

func (c *Client) ExpirationCompletedConsumer(ctx context.Context, errHandler streaming.ConsumeErrorHandler) (
	streaming.ExpirationCompletedConsumer,
	error,
) {
	return expirationCompleted.jsonConsumer(ctx, c.conn, errHandler)
}

func (c *Client) ExpirationCompletedPublisher(ctx context.Context) (
	streaming.ExpirationCompletedPublisher,
	error,
) {
	return expirationCompleted.publisher(ctx, c.conn)
}
