package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

func (c *Client) TicketCreatedPublisher(ctx context.Context) (
	streaming.TicketCreatedPublisher,
	error,
) {
	return createStreamIfNotExist[streaming.TicketCreatedMessage](ctx, c.conn, streaming.TicketCreatedStreamConfig)
}

func (c *Client) TicketCreatedConsumer(
	ctx context.Context,
	errHandler streaming.ConsumeErrorHandler,
	filterSubjects ...string,
) (
	streaming.TicketCreatedConsumer,
	error,
) {
	return consumer[streaming.TicketCreatedMessage](
		ctx,
		c,
		streaming.TicketCreatedStreamConfig,
		errHandler,
		filterSubjects...,
	)
}

func (c *Client) TicketUpdatedPublisher(ctx context.Context) (
	streaming.TicketUpdatedPublisher,
	error,
) {
	return createStreamIfNotExist[streaming.TicketUpdatedMessage](
		ctx,
		c.conn,
		streaming.TicketCreatedStreamConfig,
	)
}

func (c *Client) TicketUpdatedConsumer(
	ctx context.Context,
	errHandler streaming.ConsumeErrorHandler,
	filterSubjects ...string,
) (
	streaming.TicketUpdateConsumer,
	error,
) {
	return consumer[streaming.TicketUpdatedMessage](
		ctx,
		c,
		streaming.TicketUpdatedStreamConfig,
		errHandler,
		filterSubjects...,
	)
}
