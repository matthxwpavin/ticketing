package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/matthxwpavin/ticketing/jwtclaims"
)

const (
	ErrTypeNameInvalidParameter = "invalid_parameter"
	ErrTypeNameServiceFailure   = "service_failure"
)

type CustomError struct {
	Type string `json:"type"`
	Msg  string `json:"message"`
	Val  any    `json:"value,omitempty"`
}

func (s *CustomError) Error() string {
	return s.Msg
}

func NewInvalidParameterError(msg string, err ValidationErrors) *CustomError {
	return &CustomError{
		Type: "invalid_parameter",
		Msg:  msg,
		Val:  err,
	}
}

func NewServiceFailureError(msg, code string) *CustomError {
	return &CustomError{
		Type: "service_failure",
		Msg:  msg,
		Val:  code,
	}
}

func NewCustomErrorFrom(r *http.Response) (*CustomError, error) {
	var ce CustomError
	if err := json.NewDecoder(r.Body).Decode(&ce); err != nil {
		return nil, err
	}
	return &ce, nil
}

type UnauthorizedError struct {
	CustomError
}

var ErrUnauthorized = &UnauthorizedError{CustomError{
	Type: "unauthorized",
	Msg:  "Unauthorized",
}}

func IsAuthorized(ctx context.Context) (*jwtclaims.CustomClaims, error) {
	if claims := jwtclaims.FromContext(ctx); claims != nil {
		return claims, nil
	}
	return nil, ErrUnauthorized
}

func ValidateStruct(a any) error {
	if err := v.Struct(a); err != nil {
		return handleValidateError(err)
	}
	return nil
}

func handleValidateError(err error) error {
	switch v := err.(type) {
	case ValidationErrors:
		return NewInvalidParameterError("Invalid Parameters", v)
	default:
		return NewServiceFailureError("Validation Failed", "")
	}
}
