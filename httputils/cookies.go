package httputils

import (
	"net/http"
)

const JWTCookieName = "ticketing-jwt"

func SetJWTCookie(w http.ResponseWriter, jwt string) {
	cookie := &http.Cookie{
		Name:     JWTCookieName,
		Value:    jwt,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func DeleteJWTCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     JWTCookieName,
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func JWTCookie(r *http.Request) (*http.Cookie, error) {
	return r.Cookie(JWTCookieName)
}
