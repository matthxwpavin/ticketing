package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

func (c *Client) ChargeCreatedPublisher(ctx context.Context) (
	streaming.ChargeCreatedPublisher,
	error,
) {
	return createStreamIfNotExist[streaming.ChargeCreatedMessage](
		ctx,
		c.conn,
		streaming.ChargeCreatedStreamConfig,
	)
}

func (c *Client) ChargeCreatedConsumer(
	ctx context.Context,
	errHandler streaming.ConsumeErrorHandler,
	filterSubject string,
) (
	streaming.ChargeCreatedConsumer,
	error,
) {
	return consumer[streaming.ChargeCreatedMessage](
		ctx,
		c,
		streaming.ChargeCreatedStreamConfig,
		errHandler,
		filterSubject,
	)
}
