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
	Apods       ApodService
	Users       UserService
	Permissions PermissonService
	Otps        OtpService
	Tokens      TokenService
}

func NewModels(db *sql.DB) Models {
	return Models{
		Apods:       NewApodModel(db),
		Users:       NewUserModel(db),
		Permissions: NewPermissonModel(db),
		Otps:        NewOtpModel(db),
		Tokens:      NewTokenModel(db),
	}
}
