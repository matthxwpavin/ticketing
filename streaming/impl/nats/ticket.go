package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

func (c *Client) ticketCreatedSubject() *subject[streaming.TicketCreatedMessage] {
	return &subject[streaming.TicketCreatedMessage]{
		names:           []string{"ticket:created"},
		streamName:      "ticket:created",
		consumerName:    c.ConsumerName,
		consumerSubject: c.ConsumerSubject,
	}
}

func (c *Client) ticketUpdatedSubject() *subject[streaming.TicketUpdatedMessage] {
	return &subject[streaming.TicketUpdatedMessage]{
		names:           []string{"ticket:updated"},
		streamName:      "ticket:updated",
		consumerName:    c.ConsumerName,
		consumerSubject: c.ConsumerSubject,
	}
}

func (c *Client) TicketCreatedPublisher(ctx context.Context) (
	streaming.TicketCreatedPublisher,
	error,
) {
	return c.ticketCreatedSubject().publisher(ctx, c.conn)
}

func (c *Client) TicketCreatedConsumer(ctx context.Context, errHandler streaming.ConsumeErrorHandler) (
	streaming.TicketCreatedConsumer,
	error,
) {
	return c.ticketCreatedSubject().jsonConsumer(ctx, c.conn, errHandler)
}

func (c *Client) TicketUpdatedPublisher(ctx context.Context) (
	streaming.TicketUpdatedPublisher,
	error,
) {
	return c.ticketUpdatedSubject().publisher(ctx, c.conn)
}

func (c *Client) TicketUpdatedConsumer(ctx context.Context, errHandler streaming.ConsumeErrorHandler) (
	streaming.TicketUpdateConsumer,
	error,
) {
	return c.ticketUpdatedSubject().jsonConsumer(ctx, c.conn, errHandler)
}
