package streaming

import (
	"context"
)

type Publisher[T any] interface {
	Publish(ctx context.Context, message *T) error
}
