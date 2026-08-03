package contracts

import (
	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
)

type UserContract interface {
	CreateUser(user *domain.User) (*domain.User, error)
	FindUserByUUID(uuid string) (*domain.User, error)
	FindUserByEmailAndWebsite(email string, websiteUUID string) (*domain.User, error)
	UserExists(email string, websiteUUID string) (bool, error)
	GetUsersFromWebsite(websiteUUID string) ([]*domain.User, error)
	GetUsers() ([]*domain.User, error)
	UpdateUserByUUID(uuid string) error
	DeleteUserByUUID(uuid string) error
	DeleteUsersByUUIDS(uuid []string) error
}
