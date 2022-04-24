package main

import "net/http"

type HealthchekHandler interface {
	Check(w http.ResponseWriter, r *http.Request)
}

type healhcheckApi struct {
	app *application
}

func NewHealthcheckApi(app *application) HealthchekHandler{
	return &healhcheckApi{app: app}
}

func (h *healhcheckApi) Check(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": h.app.config.env,
			"version":     version,
		},
	}

	err := h.app.writeJSON(w, http.StatusOK, 0, env, nil)
	if err != nil {
		h.app.serverErrorResponse(w, r, err)
	}

}
