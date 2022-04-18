package mock

import (
	"javlonrahimov/apod/internal/data"
	"time"
)

var tokens = make([]data.Token, 10)

type tokenModelMock struct{}

func NewTokenMock() data.TokenService {
	return &tokenModelMock{}
}

func (m tokenModelMock) New(userId int64, ttl time.Duration, scope string) (*data.Token, error) {
	token, err := data.GenerateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}
	m.DeleteAllForUser(scope, userId)
	err = m.Insert(token)
	return token, err
}

func (m tokenModelMock) Insert(token *data.Token) error {

	tokens = append(tokens, *token)

	return nil
}

func (m tokenModelMock) DeleteAllForUser(scope string, userId int64) error {
	for i := 0; i < len(tokens); i++ {
		if tokens[i].UserId == userId {
			tokens = removeToken(tokens, i)
		}
	}
	return nil
}

func removeToken(slice []data.Token, s int) []data.Token {
	return append(slice[:s], slice[s+1:]...)
}
