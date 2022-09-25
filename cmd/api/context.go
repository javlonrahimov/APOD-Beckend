package main

import (
	"context"
	"net/http"

	"apod.api.javlonrahimov1212/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

func (a *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (a *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing context value in request context")
	}
	return user
}
