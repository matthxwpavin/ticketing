package nats

import (
	"context"
	"encoding/json"

	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/streaming"
	"github.com/matthxwpavin/ticketing/streaming/ticket"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Client struct {
	URL  string // Not required to be a property due to there is no sub-type of the Client.
	Name string // Not required to be a property due to there is no sub-type of the Client.
	conn *nats.Conn
}

// Not required to be a property due to there is no sub-type of the Client.
func NewFromEnv(ctx context.Context) *Client {
	return &Client{URL: env.NatsURL(), Name: env.NatsConnectionName()}
}

type ConnectOption func(*Client)

func WithConnectionName(name string) ConnectOption {
	return func(c *Client) { c.Name = name }
}

func Connect(ctx context.Context, url string, opts ...ConnectOption) (*Client, error) {
	return connect(ctx, &Client{URL: url}, opts...)
}

func ConnectFromEnv(ctx context.Context, opts ...ConnectOption) (*Client, error) {
	return connect(ctx, NewFromEnv(ctx), opts...)
}

func connect(ctx context.Context, c *Client, opts ...ConnectOption) (*Client, error) {
	for _, opt := range opts {
		opt(c)
	}
	return c, c.Connect(ctx)
}

// Not required to be a property due to there is no sub-type of the Client.
func (c *Client) Connect(ctx context.Context) error {
	logger := sugar.FromContext(ctx)

	var opts []nats.Option
	if c.Name != "" {
		opts = append(opts, nats.Name(c.Name))
	}

	var err error
	c.conn, err = nats.Connect(c.URL, opts...)
	if err != nil {
		logger.Errorw("NATS failed to connect", "error", err)
		return err
	}
	return nil
}

func (c *Client) Disconenct(ctx context.Context) error {
	logger := sugar.FromContext(ctx)
	if err := c.conn.Drain(); err != nil {
		logger.Errorw("unable to drain", "error", err)
	}
	c.conn.Close()
	logger.Info("NATS connection closed")
	return nil
}

func (c *Client) TicketCreatedPublisher(ctx context.Context) (
	ticket.CreatedPublisher,
	error,
) {
	return publisher[ticket.CreatedMessage](ctx, c.conn, "ticket:created", "ticket:created")
}

func (c *Client) TicketUpdatedPublisher(ctx context.Context) (
	ticket.UpdatedPublisher,
	error,
) {
	return publisher[ticket.UpdatedMessage](ctx, c.conn, "ticket:updated", "ticket:updated")
}

func publisher[T any](
	ctx context.Context,
	conn *nats.Conn,
	name string,
	subject string,
) (streaming.Publisher[T], error) {
	logger := sugar.FromContext(ctx)

	js, err := jetstream.New(conn)
	if err != nil {
		logger.Errorw("unable to new jetstream", "error", err)
		return nil, err
	}
	_, err = js.Stream(ctx, name)
	notfound := err == jetstream.ErrStreamNotFound

	if err != nil && !notfound {
		logger.Errorw("unable to get stream", "error", err, "name", name)
		return nil, err
	}
	if notfound {
		_, err := js.CreateStream(ctx, jetstream.StreamConfig{
			Name:      name,
			Subjects:  []string{subject},
			Retention: jetstream.WorkQueuePolicy,
		})
		if err != nil {
			logger.Errorw("could not create stream", "error", err, "name", name)
			return nil, err
		}
	}

	return &jetStream[T]{
		name:    name,
		subject: subject,
		js:      js,
	}, nil
}

type jetStream[T any] struct {
	name    string
	subject string
	js      jetstream.JetStream
}

func (js *jetStream[T]) Publish(ctx context.Context, message *T) error {
	logger := sugar.FromContext(ctx)
	payload, err := json.Marshal(message)
	if err != nil {
		logger.Errorw("unable to marshal", "error", err)
		return err
	}
	if _, err := js.js.Publish(ctx, js.subject, payload); err != nil {
		logger.Errorw("unable to publish", "error", err, "subject", js.subject, "paylaod", string(payload))
		return err
	}
	return nil
}
