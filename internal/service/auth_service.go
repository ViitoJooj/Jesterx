package service

import (
	"errors"
	"strings"
	"time"

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

func (s *AuthService) Register(first_name string, last_name string, email string, password string) (*domain.User, error) {
	if email == "" || len(email) > 250 || len(email) < 5 || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return nil, errors.New("Invalid email.")
	}
	if password == "" || len(password) > 50 || len(password) < 6 {
		return nil, errors.New("Invalid password.")
	}
	if first_name == "" || last_name == "" || len(first_name) > 50 || len(last_name) > 50 || len(first_name) < 2 || len(last_name) < 2 {
		return nil, errors.New("Invalid name.")
	}
	existing, err := s.userRepo.FindUserByEmail(email)
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

	if err := s.userRepo.UserRegister(*user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email string, password string) (*domain.User, error) {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("Invalid.")
	}
	if !security.CheckPassword(password, user.Password) {
		return nil, errors.New("Invalid.")
	}

	return user, nil
}

func (s *AuthService) Refresh(refreshToken string) (string, error) {
	refreshClaims, err := security.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindUserByID(refreshClaims.Sub)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	accessClaims := security.AccessTokenClaims{
		Iss:  "https://jesterx.com.br",
		Sub:  user.Id,
		Aud:  "https://api.jesterx.com.br",
		Exp:  time.Now().Add(15 * time.Minute).Unix(),
		Role: user.Role,
	}

	accessToken, err := security.AccessToken(accessClaims)
	if err != nil {
		return "", errors.New("failed to generate access token")
	}

	return accessToken, nil
}

func (s *AuthService) Me(accessToken string) (*domain.User, error) {
	claims, err := security.ParseAccessToken(accessToken)
	if err != nil {
		return nil, errors.New("invalid access token")
	}

	user, err := s.userRepo.FindUserByID(claims.Sub)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
