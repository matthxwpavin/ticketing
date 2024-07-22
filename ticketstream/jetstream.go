package ticketstream

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type JetStream struct {
	nc *nats.Conn
	js jetstream.JetStream
}

func Default() (*JetStream, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	return &JetStream{nc: nc, js: js}, nil
}

func (s *JetStream) DrainAndClose() error {
	if err := s.nc.Drain(); err != nil {
		return err
	}
	s.nc.Close()
	return nil
}

func (s *JetStream) GetStreamOrCreate(ctx context.Context, name string, subjects []string) (jetstream.Stream, error) {
	stream, err := s.js.Stream(context.Background(), name)
	if err == nil {
		return stream, nil
	}
	if err == jetstream.ErrStreamNotFound {
		stream, err := s.js.CreateStream(ctx, jetstream.StreamConfig{
			Name:      name,
			Subjects:  subjects,
			Retention: jetstream.WorkQueuePolicy,
		})
		if err != nil {
			return nil, err
		}
		return stream, nil
	}
	return nil, err
}

func (s *JetStream) GetConsumerOrCreate(ctx context.Context, name string, stream string) (jetstream.Consumer, error) {
	cmr, err := s.js.Consumer(ctx, stream, name)
	if err == nil {
		return cmr, nil
	}
	if err == jetstream.ErrConsumerNotFound {
		cmr, err := s.js.CreateConsumer(
			ctx,
			stream,
			jetstream.ConsumerConfig{
				Name:          name,
				AckPolicy:     jetstream.AckExplicitPolicy,
				DeliverPolicy: jetstream.DeliverAllPolicy,
			},
		)
		if err != nil {
			return nil, err
		}
		return cmr, nil
	}
	return nil, err
}

func (s *JetStream) PublishJSON(ctx context.Context, subject string, message any) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = s.js.Publish(ctx, subject, data)
	return err
}

func (s *JetStream) JetStream() jetstream.JetStream {
	return s.js
}

type Stream interface {
	jetstream.Stream
}

type Consumer interface {
	jetstream.Consumer
}
