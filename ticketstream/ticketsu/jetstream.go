package ticketsu

import (
	"context"

	"github.com/matthxwpavin/ticketing/ticketstream"
	"github.com/nats-io/nats.go"
)

const stream = "tickets:updated"

const subject = "tickets:updated"

const consumer = "ticketsu"

type JetStream struct {
	*ticketstream.JetStream
}

func NewJetStream(connectUrl string) (*JetStream, error) {
	js, err := ticketstream.Connect(connectUrl)
	if err != nil {
		return nil, err
	}
	return &JetStream{JetStream: js}, nil
}

func From(nc *nats.Conn) (*JetStream, error) {
	js, err := ticketstream.From(nc)
	if err != nil {
		return nil, err
	}
	return &JetStream{JetStream: js}, nil
}

func (s *JetStream) GetStreamOrCreate(ctx context.Context) error {
	_, err := s.JetStream.GetStreamOrCreate(ctx, stream, []string{subject})
	return err
}

func (s *JetStream) GetConsumerOrCreate(ctx context.Context) (ticketstream.Consumer, error) {
	return s.JetStream.GetConsumerOrCreate(ctx, consumer, stream)
}

type Message struct {
	TicketID      string  `json:"ticketID"`
	TicketTitle   string  `json:"ticketTitle"`
	TicketPrice   float64 `json:"ticketPrice"`
	UserID        string  `json:"userID"`
	TicketVersion int     `json:"ticketVersion"`
}

func (s *JetStream) Publish(ctx context.Context, message *Message) error {
	return s.JetStream.PublishJSON(ctx, subject, message)
}
