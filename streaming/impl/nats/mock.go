package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

type MockClient struct {
	streaming.AckTopicsMessages
}

func (c *MockClient) TicketCreatedPublisher(context.Context) (
	streaming.TicketCreatedPublisher,
	error,
) {
	return &mockJetStream[streaming.TicketCreatedMessage]{
		msg: c.TicketCreatedMsg,
	}, nil
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
	return &mockJetStream[streaming.TicketCreatedMessage]{
		msg: c.TicketCreatedMsg,
	}, nil
}

func (c *MockClient) TicketUpdatedConsumer(context.Context, streaming.ConsumeErrorHandler) (
	streaming.TicketUpdateConsumer,
	error,
) {
	return &mockJetStream[streaming.TicketUpdatedMessage]{
		msg: c.TicketUpdatedMsg,
	}, nil
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

func (c *MockClient) OrderCreatedConsumer(context.Context, streaming.ConsumeErrorHandler) (
	streaming.OrderCreatedConsumer,
	error,
) {
	return &mockJetStream[streaming.OrderCreatedMessage]{
		msg: c.OrderCreatedMsg,
	}, nil
}

func (c *MockClient) OrderCancelledConsumer(context.Context, streaming.ConsumeErrorHandler) (
	streaming.OrderCancelledConsumer,
	error,
) {
	return &mockJetStream[streaming.OrderCancelledMessage]{
		msg: c.OrderCancelledMsg,
	}, nil
}

type mockJetStream[T any] struct {
	msg *streaming.AcknowledgeMessage[T]
}

func (mjs *mockJetStream[T]) Publish(context.Context, *T) error { return nil }

func (mjs *mockJetStream[T]) Consume(ctx context.Context, handler streaming.JsonMessageHandler[T]) (streaming.Unsubscriber, error) {
	if mjs.msg != nil {
		handler(mjs.msg.Message, func() error { mjs.msg.Ack(); return nil })
	}
	return nil, nil
}
