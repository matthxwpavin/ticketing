package nats

import (
	"context"
	"encoding/json"

	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/streaming"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Client struct {
	URL          string // Not required to be a property due to there is no sub-type of the Client.
	Name         string // Not required to be a property due to there is no sub-type of the Client.
	ConsumerName string
	conn         *nats.Conn
}

// Not required to be a property due to there is no sub-type of the Client.
func NewFromEnv(ctx context.Context, consumerName string) *Client {
	return &Client{URL: env.NatsURL(), Name: env.NatsConnectionName(), ConsumerName: consumerName}
}

type ConnectOption func(*Client)

func WithConnectionName(name string) ConnectOption {
	return func(c *Client) { c.Name = name }
}

func Connect(ctx context.Context, url string, name string, consumerName string, opts ...ConnectOption) (*Client, error) {
	return connect(ctx, &Client{
		URL:          url,
		ConsumerName: consumerName,
	}, opts...)
}

func ConnectFromEnv(ctx context.Context, consumerName string, opts ...ConnectOption) (*Client, error) {
	return connect(ctx, NewFromEnv(ctx, consumerName), opts...)
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

type subject[T any] struct {
	name         string
	streamName   string
	consumerName string
}

func (s *subject[T]) publisher(ctx context.Context, conn *nats.Conn) (streaming.Publisher[T], error) {
	logger := sugar.FromContext(ctx)
	js, err := createStreamIfNotExist(ctx, conn, s.streamName, s.name)
	if err != nil {
		logger.Errorw("could not create stream", "error", err)
		return nil, err
	}
	return &jetStream[T]{
		name:         s.streamName,
		subject:      s.name,
		consumerName: s.consumerName,
		js:           js,
	}, nil
}

func (s *subject[T]) jsonConsumer(
	ctx context.Context,
	conn *nats.Conn,
	errHandler streaming.ConsumeErrorHandler,
) (streaming.JsonConsumer[T], error) {
	cmr, err := s.createConsumerIfNotExist(ctx, conn)
	if err != nil {
		return nil, err
	}
	return &jsonConsumer[T]{
		Consumer:   cmr,
		errHandler: errHandler,
	}, nil
}

func (s *subject[_]) createConsumerIfNotExist(
	ctx context.Context,
	conn *nats.Conn,
) (jetstream.Consumer, error) {
	logger := sugar.FromContext(ctx)
	js, err := createStreamIfNotExist(ctx, conn, s.streamName, s.name)
	if err != nil {
		logger.Errorw("could not create stream", "error", err)
		return nil, err
	}
	logger = logger.With("stream_name", s.streamName, "subject", s.name)
	cmr, err := js.Consumer(ctx, s.streamName, s.consumerName)
	notfound := err == jetstream.ErrConsumerNotFound
	if err != nil && !notfound {
		logger.Errorw("could not get a consumer", "error", err)
		return nil, err
	}
	if notfound {
		cmr, err = js.CreateConsumer(
			ctx,
			s.streamName,
			jetstream.ConsumerConfig{
				Durable:       s.consumerName,
				AckPolicy:     jetstream.AckExplicitPolicy,
				DeliverPolicy: jetstream.DeliverAllPolicy,
			},
		)
		if err != nil {
			return nil, err
		}
	}
	return cmr, nil
}

type jsonConsumer[T any] struct {
	errHandler streaming.ConsumeErrorHandler
	jetstream.Consumer
}

func (c *jsonConsumer[T]) Consume(ctx context.Context, handler streaming.JsonMessageHandler[T]) (streaming.Unsubscriber, error) {
	return consume(ctx, c.Consumer, c.errHandler, func(msg jetstream.Msg) {
		logger := sugar.FromContext(ctx)
		jsonMsg := new(T)
		if err := json.Unmarshal(msg.Data(), jsonMsg); err != nil {
			logger.Errorw("could not unmarshal a message", "error", err)
			return
		}
		handler(jsonMsg, msg.Ack)
	})
}

func consume(
	_ context.Context,
	consumer jetstream.Consumer,
	errHandler streaming.ConsumeErrorHandler,
	hanlder jetstream.MessageHandler,
) (streaming.Unsubscriber, error) {
	cctx, err := consumer.Consume(hanlder, jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
		errHandler(func() {
			consumeCtx.Drain()
			consumeCtx.Stop()
		}, err)
	}))
	return func() { cctx.Drain(); cctx.Stop() }, err
}

func createStreamIfNotExist(ctx context.Context, conn *nats.Conn, name string, subject string) (jetstream.JetStream, error) {
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
	return js, nil
}

type jetStream[T any] struct {
	name         string
	subject      string
	consumerName string
	js           jetstream.JetStream
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
