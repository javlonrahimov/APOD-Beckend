package main

import (
	"javlonrahimov/apod/internal/data"
	"javlonrahimov/apod/internal/validator"
	"net/http"
)

type ApodHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	GetByDate(w http.ResponseWriter, r *http.Request)
}

type apodApi struct {
	app *application
}

func NewApodApi(app *application) ApodHandler {
	return &apodApi{app: app}
}

func (a *apodApi) GetAll(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Q string
		data.Filters
	}
	v := validator.New()

	qs := r.URL.Query()

	input.Q = a.app.readString(qs, "q", "")

	input.Filters.Page = a.app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = a.app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = a.app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "date", "-id", "-title", "-date"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		a.app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, metadata, err := a.app.models.Apods.GetAll(input.Q, input.Filters)
	if err != nil {
		a.app.serverErrorResponse(w, r, err)
		return
	}

	err = a.app.writeJSON(w, http.StatusOK, 0, envelope{"movies": movies, "metadata": metadata}, nil)
	if err != nil {
		a.app.serverErrorResponse(w, r, err)
	}
}

func (a *apodApi) GetById(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Id int64
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Id = int64(a.app.readInt(qs, "id", 0, v))

	user, err := a.app.contextGetUser(r)

	if err != nil {
		a.app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	apod, err := a.app.models.Apods.GetById(input.Id, user.ID) // todo fetch user
	if err != nil {
		switch err {
		case data.ErrRecordNotFound:
			a.app.notFoundResponse(w, r)
		default:
			a.app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.app.writeJSON(w, http.StatusOK, 0, envelope{"apod": apod}, nil)
	if err != nil {
		a.app.serverErrorResponse(w, r, err)
	}

}

func (a *apodApi) GetByDate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Date string
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Date = a.app.readString(qs, "date", "")

	data.ValidateDate(v, input.Date)

	user, err := a.app.contextGetUser(r)

	if err != nil {
		a.app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	apod, err := a.app.models.Apods.GetByDate(input.Date, user.ID)
	if err != nil {
		switch err {
		case data.ErrRecordNotFound:
			a.app.notFoundResponse(w, r)
		default:
			a.app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.app.writeJSON(w, http.StatusOK, 0, envelope{"apod": apod}, nil)
	if err != nil {
		a.app.serverErrorResponse(w, r, err)
	}
}
