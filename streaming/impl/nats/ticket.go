package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

var (
	ticketCreated = &subject[streaming.TicketCreatedMessage]{
		name:       "ticket:created",
		streamName: "ticket:created",
	}
	ticketUpdated = &subject[streaming.TicketUpdatedMessage]{
		name:       "ticket:updated",
		streamName: "ticket:updated",
	}
)

func (c *Client) TicketCreatedPublisher(ctx context.Context) (
	streaming.TicketCreatedPublisher,
	error,
) {
	return ticketCreated.publisher(ctx, c.conn)
}

func (c *Client) TicketCreatedConsumer(ctx context.Context, errHandler streaming.ConsumeErrorHandler) (
	streaming.TicketCreatedConsumer,
	error,
) {
	return ticketCreated.jsonConsumer(ctx, c.conn, errHandler)
}

func (c *Client) TicketUpdatedPublisher(ctx context.Context) (
	streaming.TicketUpdatedPublisher,
	error,
) {
	return ticketUpdated.publisher(ctx, c.conn)
}

func (c *Client) TicketUpdatedConsumer(ctx context.Context, errHandler streaming.ConsumeErrorHandler) (
	streaming.TicketUpdateConsumer,
	error,
) {
	return ticketUpdated.jsonConsumer(ctx, c.conn, errHandler)
}
