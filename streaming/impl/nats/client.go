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
	return &Client{
		URL:          env.NatsURL(),
		Name:         env.NatsConnectionName(),
		ConsumerName: consumerName,
	}
}

type ConnectOption func(*Client)

func WithConnectionName(name string) ConnectOption {
	return func(c *Client) { c.Name = name }
}

func Connect(
	ctx context.Context,
	url string,
	name string,
	consumerName string,
	opts ...ConnectOption,
) (*Client, error) {
	return connect(ctx, &Client{
		URL:          url,
		ConsumerName: consumerName,
	}, opts...)
}

func ConnectFromEnv(
	ctx context.Context,
	consumerName string,
	opts ...ConnectOption,
) (*Client, error) {
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

func createStreamIfNotExist[T any](ctx context.Context, conn *nats.Conn, streamConfig *streaming.StreamConfig) (*jetStream[T], error) {
	logger := sugar.FromContext(ctx)

	js, err := jetstream.New(conn)
	if err != nil {
		logger.Errorw("unable to new jetstream", "error", err)
		return nil, err
	}
	_, err = js.Stream(ctx, streamConfig.Name)
	notfound := err == jetstream.ErrStreamNotFound

	if err != nil && !notfound {
		logger.Errorw("unable to get stream", "error", err, "name", streamConfig.Name)
		return nil, err
	}
	if notfound {
		_, err := js.CreateStream(ctx, jetstream.StreamConfig{
			Name:      streamConfig.Name,
			Subjects:  streamConfig.Subjects,
			Retention: jetstream.WorkQueuePolicy,
		})
		if err != nil {
			logger.Errorw("could not create stream", "error", err, "name", streamConfig.Name)
			return nil, err
		}
	}
	return &jetStream[T]{
		js:       js,
		name:     streamConfig.Name,
		subjects: streamConfig.Subjects,
	}, nil
}

type jetStream[T any] struct {
	name     string
	subjects []string
	js       jetstream.JetStream
}

func (js *jetStream[T]) consumer(
	ctx context.Context,
	consumerName string,
	errHandler streaming.ConsumeErrorHandler,
	filterSubjects ...string,
) (*jsonConsumer[T], error) {
	logger := sugar.FromContext(ctx)
	cmr, err := js.js.Consumer(ctx, js.name, consumerName)
	notfound := err == jetstream.ErrConsumerNotFound
	if err != nil && !notfound {
		logger.Errorw("could not get a consumer", "error", err)
		return nil, err
	}
	if notfound {
		cmr, err = js.js.CreateConsumer(
			ctx,
			js.name,
			jetstream.ConsumerConfig{
				Durable:        consumerName,
				AckPolicy:      jetstream.AckExplicitPolicy,
				DeliverPolicy:  jetstream.DeliverAllPolicy,
				FilterSubjects: filterSubjects,
			},
		)
		if err != nil {
			return nil, err
		}
	}
	return &jsonConsumer[T]{
		Consumer:   cmr,
		errHandler: errHandler,
	}, nil
}

func (js *jetStream[T]) Publish(ctx context.Context, message *T) error {
	logger := sugar.FromContext(ctx)
	payload, err := json.Marshal(message)
	if err != nil {
		logger.Errorw("unable to marshal", "error", err)
		return err
	}
	for _, subject := range js.subjects {
		if _, err := js.js.Publish(ctx, subject, payload); err != nil {
			logger.Errorw("unable to publish", "error", err, "subject", subject, "paylaod", string(payload))
			return err
		}
	}

	return nil
}
