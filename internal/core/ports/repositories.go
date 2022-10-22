package ports

import (
	"time"

	"apod.api.javlonrahimov1212/internal/core/domain"
)

type ApodRepository interface {
	Insert(apod *domain.Apod) error
	Get(id int64) (*domain.Apod, error)
	Update(apod *domain.Apod) error
	Delete(id int64) error
	GetAll(title, copyright, date string, filters domain.Filters) ([]*domain.Apod, domain.Metadata, error)
}

type UserRepository interface {
	Insert(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	GetForToken(tokenScope, tokenPlainText string) (*domain.User, error)
}

type PermissionRepository interface {
	GetAllForUser(userID int64) (domain.Permissions, error)
	AddForUser(userID int64, codes ...string) error
}

type TokenRepository interface {
	Insert(token *domain.Token) error
	DeleteAllForUser(scope string, userID int64) error
}
