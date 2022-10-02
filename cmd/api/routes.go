package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(a.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/health-check", a.healthchekHandler)

	router.HandlerFunc(http.MethodPost, "/v1/apods", a.requirePermission("movies:read", a.createApodHandler))
	router.HandlerFunc(http.MethodGet, "/v1/apods/:id", a.requirePermission("movies:read", a.showApodHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/apods/:id", a.requirePermission("movies:write", a.updateApodHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/apods/:id", a.requirePermission("movies:write", a.deleteApodHandler))
	router.HandlerFunc(http.MethodGet, "/v1/apods", a.requirePermission("movies:write", a.listApodsHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", a.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", a.activatUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", a.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return a.metrics(a.recoverPanic(a.enableCORS(a.rateLimit(a.authenticate(router)))))
}
