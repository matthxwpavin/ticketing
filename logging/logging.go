package logging

import (
	"context"

	"go.uber.org/zap"
)

type ctxSugarKey struct{}

var sugarKey ctxSugarKey

func WithSugar(parent context.Context, sugar *zap.SugaredLogger) context.Context {
	return context.WithValue(parent, sugarKey, sugar)
}

func FromCtxSugar(ctx context.Context) *zap.SugaredLogger {
	return ctx.Value(sugarKey).(*zap.SugaredLogger)
}

func NewSugar() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
