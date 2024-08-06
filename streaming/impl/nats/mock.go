package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

type MockClient struct {
	ticketCreatedTopic  ticketCreatedTopic
	ticketUpdatedTopic  ticketUpdatedTopic
	orderCreatedTopic   orderCreatedTopic
	orderCancelledTopic orderCancelledTopic
}

type topic[T any] struct {
	msg *streaming.AcknowledgeMessage[T]
	pub streaming.Publisher[T]
	sub streaming.JsonConsumer[T]
}

type ticketCreatedTopic struct {
	topic[streaming.TicketCreatedMessage]
}

type ticketUpdatedTopic struct {
	topic[streaming.TicketUpdatedMessage]
}

type orderCreatedTopic struct {
	topic[streaming.OrderCreatedMessage]
}

type orderCancelledTopic struct {
	topic[streaming.OrderCancelledMessage]
}

func (c *MockClient) DidTicketCreatedMessageAck() bool {
	return c.ticketCreatedTopic.msg.DidAck()
}

func (c *MockClient) DidTicketUpdatedMessageAck() bool {
	return c.ticketUpdatedTopic.msg.DidAck()
}

func (c *MockClient) DidOrderCreatedMessageAck() bool {
	return c.orderCreatedTopic.msg.DidAck()
}

func (c *MockClient) DidOrderCancelledMessageAck() bool {
	return c.orderCancelledTopic.msg.DidAck()
}

func (c *MockClient) TicketCreatedPublisher(context.Context) (
	streaming.TicketCreatedPublisher,
	error,
) {
	c.ticketCreatedTopic.pub = &mockJetStream[streaming.TicketCreatedMessage]{
		topic: &c.ticketCreatedTopic.topic,
	}
	return c.ticketCreatedTopic.pub, nil
}

func (c *MockClient) TicketUpdatedPublisher(context.Context) (
	streaming.TicketUpdatedPublisher,
	error,
) {
	c.ticketUpdatedTopic.pub = &mockJetStream[streaming.TicketUpdatedMessage]{
		topic: &c.ticketUpdatedTopic.topic,
	}
	return c.ticketUpdatedTopic.pub, nil
}

func (c *MockClient) TicketCreatedConsumer(context.Context, streaming.ConsumeErrorHandler) (
	streaming.TicketCreatedConsumer,
	error,
) {
	c.ticketCreatedTopic.sub = &mockJetStream[streaming.TicketCreatedMessage]{
		topic: &c.ticketCreatedTopic.topic,
	}
	return c.ticketCreatedTopic.sub, nil
}

func (c *MockClient) TicketUpdatedConsumer(context.Context, streaming.ConsumeErrorHandler) (
	streaming.TicketUpdateConsumer,
	error,
) {
	c.ticketUpdatedTopic.sub = &mockJetStream[streaming.TicketUpdatedMessage]{
		topic: &c.ticketUpdatedTopic.topic,
	}
	return c.ticketUpdatedTopic.sub, nil
}

func (c *MockClient) OrderCreatedPublisher(context.Context) (
	streaming.OrderCreatedPublisher,
	error,
) {
	c.orderCreatedTopic.pub = &mockJetStream[streaming.OrderCreatedMessage]{
		topic: &c.orderCreatedTopic.topic,
	}
	return c.orderCreatedTopic.pub, nil
}

func (c *MockClient) OrderCancelledPublisher(context.Context) (
	streaming.OrderCancelledPublisher,
	error,
) {
	c.orderCancelledTopic.pub = &mockJetStream[streaming.OrderCancelledMessage]{
		topic: &c.orderCancelledTopic.topic,
	}
	return c.orderCancelledTopic.pub, nil
}

func (c *MockClient) OrderCreatedConsumer(context.Context, streaming.ConsumeErrorHandler) (
	streaming.OrderCreatedConsumer,
	error,
) {
	c.orderCreatedTopic.sub = &mockJetStream[streaming.OrderCreatedMessage]{
		topic: &c.orderCreatedTopic.topic,
	}
	return c.orderCreatedTopic.sub, nil
}

func (c *MockClient) OrderCancelledConsumer(context.Context, streaming.ConsumeErrorHandler) (
	streaming.OrderCancelledConsumer,
	error,
) {
	c.orderCancelledTopic.sub = &mockJetStream[streaming.OrderCancelledMessage]{
		topic: &c.orderCancelledTopic.topic,
	}
	return c.orderCancelledTopic.sub, nil
}

type mockJetStream[T any] struct {
	topic *topic[T]
}

func (mjs *mockJetStream[T]) Publish(_ context.Context, msg *T) error {
	mjs.topic.msg = &streaming.AcknowledgeMessage[T]{Message: msg}
	return nil
}

func (mjs *mockJetStream[T]) Consume(ctx context.Context, handler streaming.JsonMessageHandler[T]) (streaming.Unsubscriber, error) {
	if mjs.topic.msg != nil {
		handler(mjs.topic.msg.Message, func() error { mjs.topic.msg.Ack(); return nil })
	}
	return nil, nil
}
