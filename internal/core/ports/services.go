package ports

import (
	"time"

	"apod.api.javlonrahimov1212/internal/core/domain"
	"apod.api.javlonrahimov1212/internal/data"
	"apod.api.javlonrahimov1212/internal/validator"
)

type ApodService interface {
	Create(title, date, explanation, hdUrl, url, mediaType, copyright string, validator *validator.Validator) (*domain.Apod, error)
	Get(id int64, validator *validator.Validator) (*domain.Apod, error)
	Update(id, title, date, explanation, hdUrl, url, mediaType string, validator *validator.Validator) (*domain.Apod, error)
	Delete(id int64, validator *validator.Validator) error
	GetPaging(title, copyright string, date time.Time, filters data.Filters, validator *validator.Validator) ([]*domain.Apod, domain.Metadata, error)
}

type UserService interface {
	Create(name, email, password string) (domain.User, error)
	Activate(tokenPlainText string) (domain.User, error)
}

type AuthService interface {
	CreateAuthToken(email, password string) string
}
