package ticketsc

import (
	"context"

	ticketstream "github.com/matthxwpavin/ticketing/jetstream"
	"github.com/matthxwpavin/ticketing/streaming/ticketsc"
	"github.com/nats-io/nats.go"
)

const stream = "tickets:created"

const subject = "tickets:created"

const consumer = "ticketsc"

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

func (s *JetStream) Publish(ctx context.Context, message *ticketsc.Message) error {
	return s.JetStream.PublishJSON(ctx, subject, message)
}
