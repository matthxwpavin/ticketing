package ticketsc

import (
	"context"

	"github.com/matthxwpavin/ticketing/ticketstream"
)

const StreamName = "tickets:created3"

const Subject = "tickets:created3"

type JetStream struct {
	*ticketstream.JetStream
}

func NewJetStream() (*JetStream, error) {
	js, err := ticketstream.Default()
	if err != nil {
		return nil, err
	}
	return &JetStream{JetStream: js}, nil
}

func (s *JetStream) GetStreamOrCreate(ctx context.Context) error {
	_, err := s.JetStream.GetStreamOrCreate(ctx, StreamName, []string{Subject})
	return err
}

func (s *JetStream) GetConsumerOrCreate(ctx context.Context) (ticketstream.Consumer, error) {
	return s.JetStream.GetConsumerOrCreate(ctx, "ticketsc", StreamName)
}

type Message struct {
	TicketID      string  `json:"ticketID"`
	TicketTitle   string  `json:"ticketTitle"`
	TicketPrice   float64 `json:"ticketPrice"`
	TicketVersion int     `json:"ticketVersion"`
}

func (s *JetStream) Publish(ctx context.Context, message *Message) error {
	return s.JetStream.PublishJSON(ctx, Subject, message)
}
