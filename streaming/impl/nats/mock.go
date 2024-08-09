package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

type MockClient struct {
	ticketCreatedTopic       ticketCreatedTopic
	ticketUpdatedTopic       ticketUpdatedTopic
	orderCreatedTopic        orderCreatedTopic
	orderCancelledTopic      orderCancelledTopic
	expirationCompletedTopic expirationCompletedTopic
}

type topic[T any] struct {
	msg chan *streaming.AcknowledgeMessage[T]
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

type expirationCompletedTopic struct {
	topic[streaming.ExpirationCompletedMessage]
}

func didAck[T any](msgCh <-chan *streaming.AcknowledgeMessage[T]) bool {
	select {
	case msg := <-msgCh:
		return msg.DidAck()
	default:
		return false
	}
}

func (c *MockClient) DidTicketCreatedMessageAck() bool {
	return didAck(c.ticketCreatedTopic.msg)
}

func (c *MockClient) DidTicketUpdatedMessageAck() bool {
	return didAck(c.ticketUpdatedTopic.msg)
}

func (c *MockClient) DidOrderCreatedMessageAck() bool {
	return didAck(c.orderCreatedTopic.msg)
}

func (c *MockClient) DidOrderCancelledMessageAck() bool {
	return didAck(c.orderCancelledTopic.msg)
}

func (c *MockClient) DidExpirationCompletedMessageAck() bool {
	return didAck(c.expirationCompletedTopic.msg)
}

func (c *MockClient) TicketCreatedPublisher(context.Context) (
	streaming.TicketCreatedPublisher,
	error,
) {
	return initTopic(&c.ticketCreatedTopic.topic).pub, nil
}

func (c *MockClient) TicketUpdatedPublisher(context.Context) (
	streaming.TicketUpdatedPublisher,
	error,
) {
	return initTopic(&c.ticketUpdatedTopic.topic).pub, nil
}

func (c *MockClient) TicketCreatedConsumer(context.Context, streaming.ConsumeErrorHandler) (
	streaming.TicketCreatedConsumer,
	error,
) {
	return initTopic(&c.ticketCreatedTopic.topic).sub, nil
}

func (c *MockClient) TicketUpdatedConsumer(context.Context, streaming.ConsumeErrorHandler) (
	streaming.TicketUpdateConsumer,
	error,
) {
	return initTopic(&c.ticketUpdatedTopic.topic).sub, nil
}

func (c *MockClient) OrderCreatedPublisher(context.Context) (
	streaming.OrderCreatedPublisher,
	error,
) {
	return initTopic(&c.orderCreatedTopic.topic).pub, nil
}

func (c *MockClient) OrderCancelledPublisher(context.Context) (
	streaming.OrderCancelledPublisher,
	error,
) {
	return initTopic(&c.orderCancelledTopic.topic).pub, nil
}

func (c *MockClient) OrderCreatedConsumer(context.Context, streaming.ConsumeErrorHandler) (
	streaming.OrderCreatedConsumer,
	error,
) {
	return initTopic(&c.orderCreatedTopic.topic).sub, nil
}

func (c *MockClient) OrderCancelledConsumer(context.Context, streaming.ConsumeErrorHandler) (
	streaming.OrderCancelledConsumer,
	error,
) {
	return initTopic(&c.orderCancelledTopic.topic).sub, nil
}

func (c *MockClient) ExpirationCompletedConsumer(context.Context, streaming.ConsumeErrorHandler) (
	streaming.ExpirationCompletedConsumer,
	error,
) {
	return initTopic(&c.expirationCompletedTopic.topic).sub, nil
}

func (c *MockClient) ExpirationCompletedPublisher(context.Context) (
	streaming.ExpirationCompletedPublisher,
	error,
) {
	return initTopic(&c.expirationCompletedTopic.topic).pub, nil
}

func initTopic[T any](topic *topic[T]) *topic[T] {
	if topic.sub == nil {
		topic.sub = &mockJetStream[T]{
			topic: topic,
		}
	}
	if topic.pub == nil {
		topic.pub = &mockJetStream[T]{
			topic: topic,
		}
	}
	if topic.msg == nil {
		topic.msg = make(chan *streaming.AcknowledgeMessage[T], 1)
	}
	return topic
}

type mockJetStream[T any] struct {
	topic *topic[T]
}

func (mjs *mockJetStream[T]) Publish(_ context.Context, msg *T) error {
	mjs.topic.msg <- &streaming.AcknowledgeMessage[T]{Message: msg}
	return nil
}

func (mjs *mockJetStream[T]) Consume(ctx context.Context, handler streaming.JsonMessageHandler[T]) (streaming.Unsubscriber, error) {
	go func() {
		msg := <-mjs.topic.msg
		handler(msg.Message, func() error { msg.Ack(); mjs.topic.msg <- msg; return nil })
	}()
	return nil, nil
}
