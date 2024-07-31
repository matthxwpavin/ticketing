package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/matthxwpavin/ticketing/jwtclaims"
)

const (
	ErrTypeNameInvalidParameter = "invalid_parameter"
	ErrTypeNameServiceFailure   = "service_failure"
)

type customError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

func (s *customError) Error() string {
	return fmt.Sprintf("type: %v, message: %v, value: %v", s.Type, s.Message, s.Value)
}

type InvalidParameterError struct {
	customError
}

func newInvalidParameterError(err ValidationErrors) *InvalidParameterError {
	return &InvalidParameterError{
		customError{
			Type:    "invalid_parameter",
			Message: "Invalid Parameters",
			Value:   err,
		}}
}

type ServiceFailureError struct {
	customError
}

func (se *ServiceFailureError) WithCode(code string) *ServiceFailureError {
	se.Value = code
	return se
}

func NewServiceFailureError(message string) *ServiceFailureError {
	return &ServiceFailureError{
		customError{
			Type:    "service_failure",
			Message: message,
		}}
}

func NewCustomErrorFrom(r *http.Response) (*customError, error) {
	var ce customError
	if err := json.NewDecoder(r.Body).Decode(&ce); err != nil {
		return nil, err
	}
	return &ce, nil
}

type UnauthorizedError struct {
	customError
}

var ErrUnauthorized = &UnauthorizedError{customError{
	Type:    "unauthorized",
	Message: "Unauthorized",
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
		return newInvalidParameterError(v)
	default:
		return NewServiceFailureError("Validation Failed")
	}
}
