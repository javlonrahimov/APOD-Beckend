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


type UserModelMock struct {}


func (m *UserModelMock) Insert(user *data.User) error {
	return nil
}

func (m *UserModelMock) GetByEmail(email string) (*data.User, error) {
	return nil, nil
}

func (m *UserModelMock) Update(user *data.User) error {
	return nil
}

func (m *UserModelMock) GetForToken(tokenScope, tokenPlaintext string) (*data.User, error) {
	return nil, nil
}