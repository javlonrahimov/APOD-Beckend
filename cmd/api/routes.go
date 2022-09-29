package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(a.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/health-check", a.healthchekHandler)

	router.HandlerFunc(http.MethodPost, "/v1/apods", a.requireAuthenticatedUser(a.createApodHandler))
	router.HandlerFunc(http.MethodGet, "/v1/apods/:id", a.requireAuthenticatedUser(a.showApodHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/apods/:id", a.requireAuthenticatedUser(a.updateApodHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/apods/:id", a.requireAuthenticatedUser(a.deleteApodHandler))
	router.HandlerFunc(http.MethodGet, "/v1/apods", a.requireAuthenticatedUser(a.listApodsHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", a.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", a.activatUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", a.createAuthenticationTokenHandler)

	return a.recoverPanic(a.rateLimit(a.authenticate(router)))
}
