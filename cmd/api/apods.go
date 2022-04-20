package main

import (
	"javlonrahimov/apod/internal/data"
	"javlonrahimov/apod/internal/validator"
	"net/http"
)

func (app *application) getAllHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Q string
		data.Filters
	}
	v := validator.New()

	qs := r.URL.Query()

	input.Q = app.readString(qs, "q", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "date", "-id", "-title", "-date"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, metadata, err := app.models.Apods.GetAll(input.Q, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, 0, envelope{"movies": movies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getById(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Id int64
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Id = int64(app.readInt(qs, "id", 0, v))

	apod, err := app.models.Apods.GetById(input.Id)
	if err != nil {
		switch err {
		case data.ErrRecordNotFound:
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, 0, envelope{"apod": apod}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) GetByDate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Date string
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Date = app.readString(qs, "date", "")

	data.ValidateDate(v, input.Date)

	apod, err := app.models.Apods.GetByDate(input.Date)
	if err != nil {
		switch err {
		case data.ErrRecordNotFound:
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, 0, envelope{"apod": apod}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
