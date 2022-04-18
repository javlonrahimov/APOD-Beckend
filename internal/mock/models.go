package mock

import (
	"javlonrahimov/apod/internal/data"
)

func NewModelsMock() data.Models {
	return data.Models{
		Apods:       NewApodsMock(),
		Users:       NewUsersMock(),
		Permissions: NewPermissionsMock(),
		Otps:        NewOtpsMock(),
		Tokens:      NewTokenMock(),
	}
}
