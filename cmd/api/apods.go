package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"apod.api.javlonrahimov1212/internal/data"
	"apod.api.javlonrahimov1212/internal/validator"
)

func (a *application) createApodHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string         `json:"title"`
		Date        string         `json:"date"`
		Explanation string         `json:"explanation"`
		HdUrl       string         `json:"hd_url"`
		Url         string         `json:"url"`
		MediaType   data.MediaType `json:"media_type"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResposne(w, r, err)
		return
	}

	apod := &data.Apod{
		Date:        time.Time{},
		Explanation: input.Explanation,
		HdUrl:       input.HdUrl,
		Url:         input.Url,
		Title:       input.Title,
		MediaType:   input.MediaType,
		CreatedAt:   time.Time{},
	}

	v := validator.New()

	if data.ValidateApod(v, apod); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.models.Apods.Insert(apod)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("v1/apods/%d", apod.ID))

	err = a.writeJSON(w, http.StatusAccepted, envelope{"apod": apod}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (a *application) showApodHandler(w http.ResponseWriter, r *http.Request) {

	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	apod, err := a.models.Apods.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"apod": apod}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) updateApodHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	apod, err := a.models.Apods.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Date        *time.Time      `json:"date"`
		Explanation *string         `json:"explanation"`
		HdUrl       *string         `json:"hd_url"`
		Url         *string         `json:"url"`
		Title       *string         `json:"title"`
		MediaType   *data.MediaType `json:"media_type"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResposne(w, r, err)
		return
	}

	if input.Date != nil {
		apod.Date = *input.Date
	}

	if input.Explanation != nil {
		apod.Explanation = *input.Explanation
	}

	if input.HdUrl != nil {
		apod.HdUrl = *input.HdUrl
	}

	if input.Url != nil {
		apod.Url = *input.Url
	}

	if input.Title != nil {
		apod.Title = *input.Title
	}

	if input.MediaType != nil {
		apod.MediaType = *input.MediaType
	}

	v := validator.New()

	if data.ValidateApod(v, apod); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.models.Apods.Update(apod)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"apod": apod}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) deleteApodHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.models.Apods.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"message": "apod successfully deleted"}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
