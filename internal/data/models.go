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
	Apods       ApodModel
	Users       UserModel
	Permissions PermissionModel
	Otps        OtpModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Apods:       ApodModel{DB: db},
		Users:       UserModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Otps:        OtpModel{DB: db},
	}
}
