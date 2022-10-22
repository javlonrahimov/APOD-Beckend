package ports

import (
	"apod.api.javlonrahimov1212/internal/core/domain"
	"apod.api.javlonrahimov1212/internal/validator"
)

type ApodService interface {
	Create(title, date, explanation, hdUrl, url, mediaType, copyright string) (*domain.Apod, *validator.Validator, error)
	Get(id int64) (*domain.Apod, *validator.Validator, error)
	Update(id int64, title, date, explanation, hdUrl, url, mediaType string) (*domain.Apod, *validator.Validator, error)
	Delete(id int64) error
	GetPaging(title, copyright, date string, filters domain.Filters, validator *validator.Validator) ([]*domain.Apod, domain.Metadata, error)
}

type UserService interface {
	Create(name, email, password string) (domain.User, error)
	Activate(tokenPlainText string) (domain.User, error)
}

type AuthService interface {
	CreateAuthToken(email, password string) string
}
