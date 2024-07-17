package pmux

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/matthxwpavin/ticketing/logging/sugar"
)

type ResponseWriter struct {
	http.ResponseWriter
	trans ut.Translator
}

func (rw *ResponseWriter) setHeaders(headers ...string) {
	var key string
	for i, h := range headers {
		if i%2 == 1 {
			rw.Header().Set(key, h)
		} else {
			key = h
		}
	}
}

func (rw *ResponseWriter) Write201() {
	rw.WriteHeader(201)
}

func (rw *ResponseWriter) Write200() {
	rw.WriteHeader(200)
}

func (rw *ResponseWriter) Writef(ctx context.Context, statusCode int, format string, v ...any) {
	data := fmt.Sprintf(format, v...)
	rw.Writeb(ctx, statusCode, []byte(data))
}

func (rw *ResponseWriter) Writeb(ctx context.Context, statusCode int, data []byte, headers ...string) {
	rw.setHeaders(headers...)
	rw.WriteHeader(statusCode)
	if _, err := rw.ResponseWriter.Write(data); err != nil {
		sugar.FromContext(ctx).Errorw("failed to write", "error", err)
	}
}

func (rw *ResponseWriter) Writeln(ctx context.Context, statusCode int, v ...any) {
	data := fmt.Sprintln(v...)
	rw.Writeb(ctx, statusCode, []byte(data))
}

func (rw *ResponseWriter) Writev(ctx context.Context, statusCode int, v ...any) {
	data := fmt.Sprint(v...)
	rw.Writeb(ctx, statusCode, []byte(data))
}

func (rw *ResponseWriter) Writee(ctx context.Context, err error) {
	switch c := err.(type) {
	case validator.ValidationErrors:
		rw.Status400().JSON(ctx, standardError{
			Type:    "invalid_parameter",
			Message: "Invalid Parameters",
			Error:   rw.requestValidationErrors(c),
		})
	case *ServiceError:
		rw.Status400().JSON(ctx, standardError{
			Type:    "service_error",
			Message: c.Message,
			Error:   c.Code,
		})
	default:
		rw.Status500().JSON(ctx, standardError{
			Type:    "unkown",
			Message: "Something went wrong",
		})
	}
}

func (rw *ResponseWriter) requestValidationErrors(errs validator.ValidationErrors) []requestValidationError {
	var res []requestValidationError
	for _, err := range errs {
		if rw.trans == nil {
			res = append(res, requestValidationError{
				Field:   err.Field(),
				Message: err.Error(),
			})
		} else {
			res = append(res, requestValidationError{
				Field:   err.Field(),
				Message: err.Translate(rw.trans),
			})
		}
	}
	return res
}

func (rw *ResponseWriter) JSON(ctx context.Context, v any) {
	rw.json(ctx, 200, v)
}

func (rw *ResponseWriter) json(ctx context.Context, statusCode int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		sugar.FromContext(ctx).Errorw("failed to marshal", "error", err)
		rw.Writev(ctx, http.StatusInternalServerError, err)
		return
	}
	headers := []string{
		"Content-Type", "application/json",
		"Content-Length", fmt.Sprint(len(data)),
	}
	rw.Writeb(ctx, statusCode, data, headers...)
}

func (rw *ResponseWriter) statusCodePrepared(statusCode int) *StatusCodePrepared {
	return &StatusCodePrepared{statusCode: statusCode, rw: rw}
}

func (rw *ResponseWriter) Status201() *StatusCodePrepared {
	return rw.statusCodePrepared(201)
}

func (rw *ResponseWriter) Status400() *StatusCodePrepared {
	return rw.statusCodePrepared(400)
}

func (rw *ResponseWriter) Status500() *StatusCodePrepared {
	return rw.statusCodePrepared(500)
}

type StatusCodePrepared struct {
	statusCode int
	rw         *ResponseWriter
}

func (c *StatusCodePrepared) JSON(ctx context.Context, v any) {
	c.rw.json(ctx, c.statusCode, v)
}

func (c *StatusCodePrepared) Write(ctx context.Context, data []byte, headers ...string) {
	c.rw.Writeb(ctx, c.statusCode, data, headers...)
}

type ServiceError struct {
	Message string
	Code    string
}

func NewServiceError(err string, code string) *ServiceError {
	return &ServiceError{Message: err, Code: code}
}

func (e *ServiceError) Error() string {
	return e.Message
}

type standardError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
}

type requestValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
