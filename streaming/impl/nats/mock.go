package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming/orderstream"
	"github.com/matthxwpavin/ticketing/streaming/ticketstream"
)

type MockClient struct{}

func (c *MockClient) TicketCreatedPublisher(ctx context.Context) (
	ticketstream.CreatedPublisher,
	error,
) {
	return &mockJetStream[ticketstream.CreatedMessage]{}, nil
}

func (c *MockClient) TicketUpdatedPublisher(ctx context.Context) (
	ticketstream.UpdatedPublisher,
	error,
) {
	return &mockJetStream[ticketstream.UpdatedMessage]{}, nil
}

func (c *MockClient) OrderCreatedPublisher(ctx context.Context) (
	orderstream.CreatedPublisher,
	error,
) {
	return &mockJetStream[orderstream.CreatedMessage]{}, nil
}

func (c *MockClient) OrderCancelledPublisher(ctx context.Context) (
	orderstream.CancelledPublisher,
	error,
) {
	return &mockJetStream[orderstream.CancelledMessage]{}, nil
}

type mockJetStream[T any] struct{}

func (mjs *mockJetStream[T]) Publish(context.Context, *T) error { return nil }
