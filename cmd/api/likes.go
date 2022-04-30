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


func (*likeApi) Like(w http.ResponseWriter, r *http.Request) {
	
}


func (*likeApi) Revert(w http.ResponseWriter, r *http.Request) {
	
}
