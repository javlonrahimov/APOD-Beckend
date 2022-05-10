package main

import "net/http"

type LikeHandler interface {
	Like(w http.ResponseWriter, r *http.Request)
	Revert(w http.ResponseWriter, r *http.Request)
}

type likeApi struct {
	app *application
}

func NewLikeApi(app *application) LikeHandler {
	return &likeApi{app: app}
}

func (l *likeApi) Like(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ApodId int64 `json:"apod_id"`
	}

	err := l.app.readJSON(w, r, &input)
	if err != nil {
		l.app.badRequestResponse(w, r, err)
		return
	}

	user, err := l.app.contextGetUser(r)

	if err != nil {
		l.app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	err = l.app.models.Likes.Insert(input.ApodId, user.ID)

	if err != nil {
		l.app.serverErrorResponse(w, r, err)
		return
	}

	err = l.app.writeJSON(w, http.StatusOK, 0, input.ApodId, nil)
	if err != nil {
		l.app.serverErrorResponse(w, r, err)
	}
}

func (l *likeApi) Revert(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ApodId int64 `json:"apod_id"`
	}

	err := l.app.readJSON(w, r, &input)
	if err != nil {
		l.app.badRequestResponse(w, r, err)
		return
	}

	user, err := l.app.contextGetUser(r)

	if err != nil {
		l.app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	err = l.app.models.Likes.Delete(input.ApodId, user.ID)

	if err != nil {
		l.app.serverErrorResponse(w, r, err)
		return
	}

	err = l.app.writeJSON(w, http.StatusOK, 0, input.ApodId, nil)
	if err != nil {
		l.app.serverErrorResponse(w, r, err)
	}
}
