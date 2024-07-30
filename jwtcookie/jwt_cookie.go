package jwtcookie

import "net/http"

const Name = "ticketing-jwt"

func New(jwt string) *http.Cookie {
	return &http.Cookie{
		Name:     Name,
		Value:    jwt,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
}

func From(r *http.Request) (*http.Cookie, error) {
	return r.Cookie(Name)
}
