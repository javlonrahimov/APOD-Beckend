package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Apods ApodModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Apods: ApodModel{DB: db},
	}
}
