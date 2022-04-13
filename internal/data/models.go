package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrOtpExpired     = errors.New("error otp expired")
)

type Models struct {
	Apods       ApodModel
	Users       UserModel
	Permissions PermissionModel
	Otps        OtpModel
	Tokens      TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Apods:       ApodModel{DB: db},
		Users:       UserModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Otps:        OtpModel{DB: db},
		Tokens:      TokenModel{DB: db},
	}
}
