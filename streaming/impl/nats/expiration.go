package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

func (c *Client) ExpirationCompletedConsumer(
	ctx context.Context,
	errHandler streaming.ConsumeErrorHandler,
	filterSubject string,
) (
	streaming.ExpirationCompletedConsumer,
	error,
) {
	return consumer[streaming.ExpirationCompletedMessage](
		ctx,
		c,
		streaming.ExpirationCompletedStreamConfig,
		errHandler,
		filterSubject,
	)
}

func (c *Client) ExpirationCompletedPublisher(ctx context.Context) (
	streaming.ExpirationCompletedPublisher,
	error,
) {
	return createStreamIfNotExist[streaming.ExpirationCompletedMessage](
		ctx,
		c.conn,
		streaming.ExpirationCompletedStreamConfig,
	)
}
