package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type Models struct {
	Apods  ApodModel
	Tokens TokenModel
	Users  UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Apods:  ApodModel{DB: db},
		Tokens: TokenModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
