package service

import (
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repository"
	"github.com/ViitoJooj/Jesterx/internal/security"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(first_name, last_name, email, password string) (*domain.User, error) {
	if email == "" || password == "" || len(email) > 250 || len(password) > 50 || len(email) < 5 || len(password) < 6 {
		return nil, errors.New("Invalid email.")
	}
	if first_name == "" || last_name == "" || len(first_name) > 50 || len(last_name) > 50 || len(first_name) < 2 || len(last_name) < 2 {
		return nil, errors.New("Invalid name.")
	}
	existing, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("Invalid.")
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := domain.NewUser(first_name, last_name, email, hashedPassword)

	if err := s.userRepo.Save(*user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("Invalid.")
	}

	if user.Banned == true {
		return nil, errors.New("Invalid.")
	}

	if !security.CheckPassword(password, user.Password) {
		return nil, errors.New("Invalid.")
	}

	return user, nil
}
