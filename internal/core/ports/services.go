package ports

import (
	"time"

	"apod.api.javlonrahimov1212/internal/core/domain"
	"apod.api.javlonrahimov1212/internal/data"
)

type ApodService interface {
	Create(title, date, explanation, hdUrl, url, mediaType string) error
	Get(id string) (domain.Apod, error)
	Update(id, title, date, explanation, hdUrl, url, mediaType string) (domain.Apod, error)
	Delete(id string) error
	GetPaging(title, copyright string, date time.Time, filters data.Filters)
}

type UserService interface {
	Create(name, email, password string) (domain.User, error)
	Activate(tokenPlainText string) (domain.User, error)
}

type AuthService interface {
	CreateAuthToken(email, password string) string
}
