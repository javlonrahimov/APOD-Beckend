package mock

import (
	"javlonrahimov/apod/internal/data"
)

type ModelsMock struct {
	Apods       data.ApodService
	Users       data.UserService
	Permissions data.PermissonService
	Otps        data.OtpService
	Tokens      data.TokenService
}

// func NewModelsMock(db *sql.DB) ModelsMock {
// 	return ModelsMock{
// 		Apods:       NewApodModel(db),
// 		Users:       &UserModelMock{},
// 		Permissions: NewPermissonModel(db),
// 		Otps:        NewOtpModel(db),
// 		Tokens:      NewTokenModel(db),
// 	}
// }
