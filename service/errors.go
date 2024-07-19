package service

import (
	"context"

	"github.com/matthxwpavin/ticketing/jwtutils"
	"github.com/matthxwpavin/ticketing/validator"
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

func NewInvalidParameterError(msg string, err validator.ValidationErrors) *CustomError {
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

type UnauthorizedError struct {
	CustomError
}

var ErrUnauthorized = &UnauthorizedError{CustomError{
	Type: "unauthorized",
	Msg:  "Unauthorized",
}}

func IsAuthorized(ctx context.Context) (*jwtutils.CustomClaims, error) {
	if claims := jwtutils.ClaimsFromContext(ctx); claims != nil {
		return claims, nil
	}
	return nil, ErrUnauthorized
}

func HandleValidateError(err error) error {
	switch v := err.(type) {
	case validator.ValidationErrors:
		return NewInvalidParameterError("Invalid Parameters", v)
	default:
		return NewServiceFailureError("Validation Failed", "")
	}
}
