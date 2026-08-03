package usecases

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repositories/contracts"
	"github.com/google/uuid"
)

type RegisterUserUseCase struct {
	db       *sql.DB
	userRepo contracts.UserContract
}

func NewRegisterUserUseCase(db *sql.DB, userRepo contracts.UserContract) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		db:       db,
		userRepo: userRepo,
	}
}

func (u *RegisterUserUseCase) Register(input *domain.User, website uuid.UUID) (*domain.User, error) {

	user, err := domain.NewUser(website.String(), "", input.Name, input.Email, "", input.Password, "", false, false, false, u.db)
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
