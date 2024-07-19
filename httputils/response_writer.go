package httputils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/service"
)

func JSONWithStatusCode(ctx context.Context, w http.ResponseWriter, statusCode int, v any) {
	logger := sugar.FromContext(ctx)
	data, err := json.Marshal(v)
	if err != nil {
		logger.Errorw("Failed to marshal", "error", err)
		w.WriteHeader(500)
		return
	}
	header := w.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Content-Length", fmt.Sprint(len(data)))
	w.WriteHeader(statusCode)
	fmt.Fprint(w, string(data))
}

func JSON(ctx context.Context, w http.ResponseWriter, v any) {
	JSONWithStatusCode(ctx, w, 200, v)
}

func JSON400(ctx context.Context, w http.ResponseWriter, v any) {
	JSONWithStatusCode(ctx, w, 400, v)
}

func JSON401(ctx context.Context, w http.ResponseWriter, v any) {
	JSONWithStatusCode(ctx, w, 401, v)
}

func JSON201(ctx context.Context, w http.ResponseWriter, v any) {
	JSONWithStatusCode(ctx, w, 201, v)
}

func Error(ctx context.Context, w http.ResponseWriter, err error) {
	switch err.(type) {
	case *service.CustomError:
		JSON400(ctx, w, err)
	case *service.UnauthorizedError:
		JSON401(ctx, w, err)
	default:
		w.WriteHeader(400)
	}
}

func DecodeJSON(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
