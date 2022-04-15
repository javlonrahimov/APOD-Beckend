package mock

import (
	"javlonrahimov/apod/internal/data"
	"time"
)

var mockUser = &data.User{
	ID:        1,
	CreatedAt: time.Time{},
	Name:      "Mock",
	Email:     "mock@gmail.com",
	Activated: false,
	Version:   1,
}

type UserModel struct {}

func (m *UserModel) Insert()
