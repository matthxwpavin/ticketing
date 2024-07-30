package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming/ticket"
)

type MockClient struct{}

func (c *MockClient) TicketCreatedPublisher(ctx context.Context) (
	ticket.CreatedPublisher,
	error,
) {
	return &mockJetStream[ticket.CreatedMessage]{}, nil
}

func (c *MockClient) TicketUpdatedPublisher(ctx context.Context) (
	ticket.UpdatedPublisher,
	error,
) {
	return &mockJetStream[ticket.UpdatedMessage]{}, nil
}

type mockJetStream[T any] struct{}

func (mjs *mockJetStream[T]) Publish(context.Context, *T) error { return nil }
