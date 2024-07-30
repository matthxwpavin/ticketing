package middleware

import (
	"net/http"

	"github.com/matthxwpavin/ticketing/jwtclaims"
	"github.com/matthxwpavin/ticketing/jwtcookie"
	"github.com/matthxwpavin/ticketing/logging/sugar"
)

func PopulateLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, err := sugar.New()
		if err != nil {
			http.Error(w, "Failed to build logger", 500)
			return
		}
		h.ServeHTTP(
			w,
			r.WithContext(sugar.WithContext(r.Context(), logger)),
		)
	})
}

// PopulateJWTClaims try to verify a JWT from Cookie, parse it
// into CustomClaims to attatch to the request's context.
//
// If any error occurred on any step, it logs of waring level then
// continue serve the request.
func PopulateJWTClaims(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			h.ServeHTTP(w, r)
		}()

		ctx := r.Context()
		logger := sugar.FromContext(ctx)

		cookie, err := jwtcookie.From(r)
		if err != nil {
			logger.Warnw("No Cookie found", "error", err)
			return
		}
		if err := cookie.Valid(); err != nil {
			logger.Warnw("Cookie is invalid", "error", err)
			return
		}
		claims, err := jwtclaims.Parse(cookie.Value)
		if err != nil {
			logger.Warnw("Failed to verify JWT", "error", err)
			return
		}
		r = r.WithContext(jwtclaims.WithContext(ctx, claims))
	})
}
