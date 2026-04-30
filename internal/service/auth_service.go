package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repository"
	"github.com/ViitoJooj/Jesterx/pkg/encryption"
	"github.com/ViitoJooj/Jesterx/pkg/jwt"
	goredis "github.com/redis/go-redis/v9"
)

type AuthService interface {
	Register(name, email, password, role, cpf string) (*domain.User, error)
	Login(email, password string) (accessToken, refreshToken string, err error)
	Refresh(refreshToken string) (accessToken string, err error)
	Logout(refreshToken string) error
}

type authService struct {
	userDomain  *domain.User
	userRepo    repository.UsersRepository
	redisClient *goredis.Client
}

func NewAuthService(userRepo repository.UsersRepository, redisClient *goredis.Client) AuthService {
	return &authService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

func (s *authService) Register(name, email, password, role, cpf string) (*domain.User, error) {
	user, err := s.userDomain.NewUser(name, email, password, role, cpf)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := encryption.Hash(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	return s.userRepo.InsertUser(user)
}

func (s *authService) Login(email, password string) (string, string, error) {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", "", err
		}
		return "", "", err
	}

	match, err := encryption.Match(password, user.Password)
	if err != nil {
		return "", "", err
	}
	if !match {
		return "", "", err
	}

	accessToken, err := jwt.GenAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.GenRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	key := fmt.Sprintf("refresh:%s", claims.ID)
	if err := s.redisClient.Set(context.Background(), key, user.Uuid, ttl).Err(); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) Refresh(refreshToken string) (string, error) {
	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("refresh:%s", claims.ID)
	userID, err := s.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.FindUserById(userID)
	if err != nil {
		return "", err
	}

	return jwt.GenAccessToken(user)
}

func (s *authService) Logout(refreshToken string) error {
	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("refresh:%s", claims.ID)
	return s.redisClient.Del(context.Background(), key).Err()
}
