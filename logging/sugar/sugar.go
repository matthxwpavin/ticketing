package sugar

import (
	"context"

	"go.uber.org/zap"
)

type ctxSugarKey struct{}

var sugarKey ctxSugarKey

func WithContext(parent context.Context, sugar *zap.SugaredLogger) context.Context {
	return context.WithValue(parent, sugarKey, sugar)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	return ctx.Value(sugarKey).(*zap.SugaredLogger)
}

func New() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
