package service

type CustomError struct {
	Type string `json:"type"`
	Msg  string `json:"message"`
	Val  any    `json:"value,omitempty"`
}

func (s *CustomError) Error() string {
	return s.Msg
}

func NewInvalidParameterError(msg string, val any) *CustomError {
	return &CustomError{
		Type: "invalid_parameter",
		Msg:  msg,
		Val:  val,
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
