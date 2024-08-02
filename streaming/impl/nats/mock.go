package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

type MockClient struct{}

func (c *MockClient) TicketCreatedPublisher(context.Context) (
	streaming.TicketCreatedPublisher,
	error,
) {
	return &mockJetStream[streaming.TicketCreatedMessage]{}, nil
}

func (c *MockClient) TicketUpdatedPublisher(context.Context) (
	streaming.TicketUpdatedPublisher,
	error,
) {
	return &mockJetStream[streaming.TicketUpdatedMessage]{}, nil
}

func (c *MockClient) TicketCreatedConsumer(context.Context, streaming.ConsumeErrorHandler) (
	streaming.TicketCreatedConsumer,
	error,
) {
	return &mockJetStream[streaming.TicketCreatedMessage]{}, nil
}

func (c *MockClient) OrderCreatedPublisher(context.Context) (
	streaming.OrderCreatedPublisher,
	error,
) {
	return &mockJetStream[streaming.OrderCreatedMessage]{}, nil
}

func (c *MockClient) OrderCancelledPublisher(context.Context) (
	streaming.OrderCancelledPublisher,
	error,
) {
	return &mockJetStream[streaming.OrderCancelledMessage]{}, nil
}

type mockJetStream[T any] struct{}

func (mjs *mockJetStream[T]) Publish(context.Context, *T) error { return nil }

func (mjs *mockJetStream[T]) Consume(context.Context, streaming.JsonMessageHandler[T]) (streaming.Unsubscriber, error) {
	return nil, nil
}
