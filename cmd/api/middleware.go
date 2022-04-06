package main

import (
	"javlonrahimov/apod/internal/data"
	// "javlonrahimov/apod/internal/validator"
	"net/http"
	"strings"
)

func (app *application) validateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// token := headerParts[1]
		// v := validator.New()

		next.ServeHTTP(w, r)

	})
}