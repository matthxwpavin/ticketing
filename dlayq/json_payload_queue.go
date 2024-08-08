package dlayq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hdt3213/delayqueue"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/redis/go-redis/v9"
)

// JsonPayloadQueue is a delayed queue assumes that messages/payload production and
// consumption in the queue must be encoded and decoded using json.Marshal/json.Unmarshal.
// When unmarshaling has an error occurred then the queue itself do a good acknowledge(return true)
// to do not want to consume the message/payload again by re-delivery mechanic.
type JsonPayloadQueue[T any] struct {
	q *delayqueue.DelayQueue
}

func NewJsonPayloadQueue[T any](
	ctx context.Context,
	name string,
	c *redis.Client,
	process func(*T) bool,
) *JsonPayloadQueue[T] {
	q := NewQueue(ctx, name, c, func(payloadStr string) bool {
		payload := new(T)
		if err := json.Unmarshal([]byte(payloadStr), payload); err != nil {
			sugar.FromContext(ctx).
				With("queue", name).
				Errorw("could not unmarshal payload", "error", err)
			// Mark the message as acknowleded to don't want to consume again
			// from re-delivery mechanic, due to the message is invalid JSON.
			return true
		}
		return process(payload)
	})
	return &JsonPayloadQueue[T]{q: q}

}

func NewQueue(ctx context.Context, name string, c *redis.Client, process func(string) bool) *delayqueue.DelayQueue {
	q := delayqueue.NewQueue(name, c, process)
	return q.WithConcurrent(4).
		WithDefaultRetryCount(3).
		WithMaxConsumeDuration(6 * time.Minute)
}
