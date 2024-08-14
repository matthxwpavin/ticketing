package jwtcookie

import (
	"net/http"

	"github.com/matthxwpavin/ticketing/env"
)

const Name = "ticketing-jwt"

func New(jwt string) *http.Cookie {
	return &http.Cookie{
		Name:     Name,
		Value:    jwt,
		Path:     "/",
		Secure:   env.DEV.Value() != "dev",
		HttpOnly: true,
	}
}

func From(r *http.Request) (*http.Cookie, error) {
	return r.Cookie(Name)
}
