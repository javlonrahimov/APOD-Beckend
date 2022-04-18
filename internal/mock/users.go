package mock

import (
	"javlonrahimov/apod/internal/data"
	"time"
)

var users = make([]data.User, 0)

type userModelMock struct{}

func NewUsersMock() data.UserService {
	return &userModelMock{}
}

func (m *userModelMock) Insert(user *data.User) error {
	user.CreatedAt = time.Now()
	user.ID = int64(len(users))
	users = append(users, *user)
	return nil
}

func (m *userModelMock) GetByEmail(email string) (*data.User, error) {
	for i := 0; i < len(users); i++ {
		if users[i].Email == email {
			return &users[i], nil
		}
	}
	return nil, data.ErrRecordNotFound
}

func (m *userModelMock) Update(user *data.User) error {

	_, err := m.GetByEmail(user.Email)
	if err == nil {
		return data.ErrDuplicateEmail
	}

	for i := 0; i < len(users); i++ {
		if users[i].ID == user.ID {
			users[i] = *user
			return nil
		}
	}

	return data.ErrEditConflict
}

func (m *userModelMock) GetForToken(tokenScope, tokenPlaintext string) (*data.User, error) {
	return nil, nil
}
