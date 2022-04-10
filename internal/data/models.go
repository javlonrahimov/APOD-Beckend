package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Apods ApodModel
	Users UserModel
	Permissions PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Apods: ApodModel{DB: db},
		Users: UserModel{DB: db},
	}
}
