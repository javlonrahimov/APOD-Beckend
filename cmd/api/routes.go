package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes(handlers *Handlers) http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", handlers.HealthCheck.Check)

	router.HandlerFunc(http.MethodPost, "/v1/users", handlers.Users.Register)
	router.HandlerFunc(http.MethodPost, "/v1/users/verify", handlers.Users.Verify)
	router.HandlerFunc(http.MethodPost, "/v1/users/login", handlers.Users.Login)

	router.HandlerFunc(http.MethodGet, "/v1/apods", handlers.Apods.GetAll)
	router.HandlerFunc(http.MethodGet, "/v1/apod-by-id", handlers.Apods.GetById)
	router.HandlerFunc(http.MethodGet, "/v1/apod-by-date", handlers.Apods.GetByDate)

	router.HandlerFunc(http.MethodPost, "/v1/like", handlers.Likes.Like)
	router.HandlerFunc(http.MethodPost, "/v1/revert-like", handlers.Likes.Revert)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
