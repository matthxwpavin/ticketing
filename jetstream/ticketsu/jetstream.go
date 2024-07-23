package ticketsu

import (
	"context"

	ticketstream "github.com/matthxwpavin/ticketing/jetstream"
	"github.com/matthxwpavin/ticketing/streaming/ticketsu"
	"github.com/nats-io/nats.go"
)

type Publisher interface {
	Publish(context.Context, *ticketsu.Message) error
}

const stream = "tickets:updated"

const subject = "tickets:updated"

const consumer = "ticketsu"

type JetStream struct {
	*ticketstream.JetStream
}

func NewJetStream(ctx context.Context, connectUrl string) (*JetStream, error) {
	nc, err := nats.Connect(connectUrl)
	if err != nil {
		return nil, err
	}
	return new(ctx, nc)
}

func From(ctx context.Context, nc *nats.Conn) (*JetStream, error) {
	return new(ctx, nc)
}

func new(ctx context.Context, nc *nats.Conn) (*JetStream, error) {
	js, err := ticketstream.From(nc)
	if err != nil {
		return nil, err
	}
	if err := js.CreateStreamIfNotExists(ctx, stream, []string{subject}); err != nil {
		return nil, err
	}
	return &JetStream{JetStream: js}, nil
}

func (s *JetStream) GetConsumerOrCreate(ctx context.Context) (ticketstream.Consumer, error) {
	return s.JetStream.GetConsumerOrCreate(ctx, consumer, stream)
}

func (s *JetStream) Publish(ctx context.Context, message *ticketsu.Message) error {
	return s.JetStream.PublishJSON(ctx, subject, message)
}
