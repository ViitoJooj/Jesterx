package repository

import "github.com/ViitoJooj/Jesterx/internal/domain"

type UserRepository interface {
	UserRegister(user domain.User) error
	FindUserByEmail(email string) (*domain.User, error)
	FindUserByEmailAndWebsite(email string, website string) (*domain.User, error)
	FindUserByID(id string) (*domain.User, error)
	DeleteUserByID(id string) error
	DeleteExpiredUnverifiedUsers() error
	UpdateVerifiedEmailToTrue(id string) error
	UpdateUserProfile(id string, data domain.UpdateProfileData) error
}
