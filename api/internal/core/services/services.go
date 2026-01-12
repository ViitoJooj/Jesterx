package services

import (
	"jesterx-core/internal/core/domain"
	"jesterx-core/internal/core/ports"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(r ports.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Create(u *domain.User) error {
	return s.repo.Create(u)
}
