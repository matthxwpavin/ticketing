package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming/ticketstream"
)

func (c *Client) TicketCreatedPublisher(ctx context.Context) (
	ticketstream.CreatedPublisher,
	error,
) {
	return publisher[ticketstream.CreatedMessage](ctx, c.conn, "ticket:created", "ticket:created")
}

func (c *Client) TicketUpdatedPublisher(ctx context.Context) (
	ticketstream.UpdatedPublisher,
	error,
) {
	return publisher[ticketstream.UpdatedMessage](ctx, c.conn, "ticket:updated", "ticket:updated")
}
