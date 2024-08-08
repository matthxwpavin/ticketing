package jsondq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hdt3213/delayqueue"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/redis/go-redis/v9"
)

// Queue is a delayed queue assumes that messages/payload production and
// consumption in the queue must be encoded and decoded using json.Marshal/json.Unmarshal.
// When unmarshaling has an error occurred then the queue itself do a good acknowledge(return true)
// to do not want to consume the message/payload again by re-delivery mechanic.
type Queue[T any] struct {
	name string
	q    *delayqueue.DelayQueue
}

func New[T any](
	ctx context.Context,
	name string,
	c *redis.Client,
	process func(*T) bool,
) *Queue[T] {
	return &Queue[T]{
		name: name,
		q: applyOptions(delayqueue.NewQueue(name, c, func(payloadStr string) bool {
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
		})),
	}
}

func (q *Queue[T]) SendImmediatelyMsg(payload *T) error {
	return q.sendJSON(payload, func(s string) error {
		return q.q.SendDelayMsg(s, 0)
	})
}

func (q *Queue[T]) SendDelayMsg(payload *T, d time.Duration) error {
	return q.sendJSON(payload, func(s string) error {
		return q.q.SendDelayMsg(s, d)
	})
}

func (q *Queue[T]) SendScheduleMsg(payload *T, t time.Time) error {
	return q.sendJSON(payload, func(s string) error {
		return q.q.SendScheduleMsg(s, t)
	})
}

func (q *Queue[T]) sendJSON(payload *T, sender func(string) error) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return sender(string(data))
}

func (q *Queue[T]) Listen(ctx context.Context) <-chan struct{} {
	go q.waitStop(ctx)
	return q.q.StartConsume()
}

func (q *Queue[T]) waitStop(ctx context.Context) {
	<-ctx.Done()
	q.q.StopConsume()
	sugar.FromContext(ctx).Infow("stopped consumption", "queue", q.name)
}

func applyOptions(q *delayqueue.DelayQueue) *delayqueue.DelayQueue {
	return q.WithConcurrent(4).
		WithDefaultRetryCount(3).
		WithMaxConsumeDuration(6 * time.Minute)
}
