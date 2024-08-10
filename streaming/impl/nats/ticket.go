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

func consumer[T any](
	ctx context.Context,
	c *Client,
	config *streaming.StreamConfig,
	errHandler streaming.ConsumeErrorHandler,
	filterSubjects ...string,
) (*jsonConsumer[T], error) {
	js, err := createStreamIfNotExist[T](ctx, c.conn, config)
	if err != nil {
		return nil, err
	}
	return js.consumer(ctx, c.ConsumerName, errHandler, filterSubjects...)
}
