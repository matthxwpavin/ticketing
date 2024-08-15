package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

func (c *Client) ChargeCreatedPublisher(ctx context.Context) (
	streaming.PaymentCreatedPublisher,
	error,
) {
	return createStreamIfNotExist[streaming.PaymentCreatedMessage](
		ctx,
		c.conn,
		streaming.PaymentCreatedStreamConfig,
	)
}

func (c *Client) PaymentCreatedConsumer(
	ctx context.Context,
	errHandler streaming.ConsumeErrorHandler,
	filterSubject string,
) (
	streaming.PaymentCreatedConsumer,
	error,
) {
	return consumer[streaming.PaymentCreatedMessage](
		ctx,
		c,
		streaming.PaymentCreatedStreamConfig,
		errHandler,
		filterSubject,
	)
}
