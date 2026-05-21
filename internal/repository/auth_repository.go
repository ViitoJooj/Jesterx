package repository

import (
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
)

type UserRepository interface {
	UserRegister(user domain.User) error
	CompanyRegister(company domain.Company) error
	FindUserByEmail(email string) (*domain.User, error)
	FindUserByEmailAndWebsite(email string, websiteID string) (*domain.User, error)
	FindUserByID(id string) (*domain.User, error)
	FindCompanyByOwnerUserID(ownerUserID string) (*domain.Company, error)
	DeleteUserByID(id string) error
	DeactivateUserByID(id string, deleteAfter time.Time) error
	DeleteExpiredUnverifiedUsers() error
	DeleteExpiredDeactivatedUsers() error
	UpdateVerifiedEmailToTrue(id string) error
	UpdateVerifiedEmailToTrueByWebsite(id string, websiteID string) error
	UpdateUserProfile(id string, data domain.UpdateProfileData) error

	ListUserAddresses(userID string) ([]*domain.UserAddress, error)
	GetDefaultUserAddress(userID string) (*domain.UserAddress, error)
	CreateUserAddress(addr domain.UserAddress) error
	UpdateUserAddress(id, userID string, data domain.UpsertAddressData) error
	DeleteUserAddress(id, userID string) error
	SetDefaultAddress(id, userID string) error
}
