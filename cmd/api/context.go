package main

import (
	"context"
	"errors"
	"javlonrahimov/apod/internal/data"
	"net/http"
)

var (
	ErrNoUserInContext = errors.New("no user in context")
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) (*data.User, error) {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		return nil, ErrNoUserInContext
	}
	return user, nil
}
