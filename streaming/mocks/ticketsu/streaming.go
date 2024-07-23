package ticketsu

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming/ticketsu"
)

type Publisher struct {
}

func (pub *Publisher) Publish(context.Context, *ticketsu.Message) error {
	return nil
}
