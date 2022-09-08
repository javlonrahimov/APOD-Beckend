package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(a.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/health-check", a.healthchekHandler)

	router.HandlerFunc(http.MethodPost, "/v1/apods", a.createApodHandler)
	router.HandlerFunc(http.MethodGet, "/v1/apods/:id", a.showApodHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/apods/:id", a.updateApodHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/apods/:id", a.deleteApodHandler)

	return router
}
