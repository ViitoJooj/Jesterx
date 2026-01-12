package ports

import (
	"jesterx-core/internal/core/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindAll() ([]*domain.User, error)
	FindByID(id string) (*domain.User, error)
	DeleteByID(id string) error
	Update(user *domain.User) error
}
