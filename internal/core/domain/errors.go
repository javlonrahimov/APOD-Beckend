package domain

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrInvaliMediaTypeFormat = errors.New("invalid media_type format")
)