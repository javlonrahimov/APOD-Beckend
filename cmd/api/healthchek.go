package main

import (
	"net/http"
)

func (a *application) healthchekHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": a.config.env,
			"version":     version,
		},
	}

	err := a.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
