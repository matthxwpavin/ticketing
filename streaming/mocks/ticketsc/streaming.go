package ticketsc

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming/ticketsc"
)

type Publisher struct {
}

func (pub *Publisher) Publish(context.Context, *ticketsc.Message) error {
	return nil
}
