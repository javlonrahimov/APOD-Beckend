package main

import "net/http"

func (app *application) healthchekHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, 0, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
