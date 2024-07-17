package pmux

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Request struct {
	r         *http.Request
	validator *validator.Validate
}

func (r *Request) Request() *http.Request {
	return r.r
}

func (r *Request) Context() context.Context {
	return r.r.Context()
}

func (r *Request) MayValidate(v any) *validation {
	return &validation{val: v, r: r.r.Body}
}

func (r *Request) Validate(v any) *validation {
	return &validation{
		val:       v,
		r:         r.r.Body,
		validator: r.validator,
		force:     true,
	}
}

func (r *Request) JSON(ctx context.Context, v any) error {
	return decodeJSON(r.r.Body, v)
}

type validation struct {
	val       any
	r         io.Reader
	validator *validator.Validate
	force     bool
}

func (v *validation) JSON() error {
	if v.force && v.validator == nil {
		return errors.New("validation is forced, but no validator")
	}
	return decodeJSON(v.r, v.val, withValidate(v.validator))
}

type decodeValidator struct {
	validate *validator.Validate
}

type decodeValidatorOptions func(dec *decodeValidator)

func withValidate(validate *validator.Validate) decodeValidatorOptions {
	return func(dec *decodeValidator) { dec.validate = validate }
}

func decodeJSON(r io.Reader, v any, options ...decodeValidatorOptions) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}
	dec := new(decodeValidator)
	for _, option := range options {
		option(dec)
	}
	if dec.validate != nil {
		return dec.validate.Struct(v)
	}
	return nil
}
