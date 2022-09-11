package main

import (
	"fmt"
	"net/http"
)

func (a *application) logError(r *http.Request, err error) {
	a.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (a *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := a.writeJSON(w, status, env, nil)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(500)
	}
}

func (a *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	a.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (a *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	a.errorResponse(w, r, http.StatusNotFound, message)
}

func (a *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	a.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (a *application) badRequestResposne(w http.ResponseWriter, r *http.Request, err error) {
	a.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (a *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	a.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (a *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	a.errorResponse(w, r, http.StatusConflict, message)
}
