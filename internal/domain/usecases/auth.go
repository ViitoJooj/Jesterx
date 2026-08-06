package usecases

import (
	"errors"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/ViitoJooj/verkoupe/internal/services"
	"github.com/google/uuid"
)

type AuthUseCase struct {
	userRepo contracts.UserContract
}

func NewAuthUseCase(userRepo contracts.UserContract) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
	}
}

func (u *AuthUseCase) Register(input *domain.User, website uuid.UUID) (*domain.User, error) {
	hashedPassword, err := services.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(website.String(), "", input.Name, input.Email, "", hashedPassword, "", false, false, false)
	if err != nil {
		return nil, err
	}

	exists, err := u.userRepo.UserExists(input.Email, website.String())
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	createdUser, err := u.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u *AuthUseCase) Login(email, password, websiteUUID string) (*domain.User, error) {
	user, err := u.userRepo.FindUserByEmailAndWebsite(email, websiteUUID)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !services.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
