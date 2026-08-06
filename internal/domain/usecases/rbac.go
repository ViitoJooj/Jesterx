package usecases

import (
	"errors"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/google/uuid"
)

type CreateRbacUseCase struct {
	rbacRepo contracts.RbacContract
}

func NewCreateRbacUseCase(rbacRepo contracts.RbacContract) *CreateRbacUseCase {
	return &CreateRbacUseCase{
		rbacRepo: rbacRepo,
	}
}

func (u *CreateRbacUseCase) Create(input *domain.Rbac, website uuid.UUID) (*domain.Rbac, error) {

	rbac, err := domain.NewRbac(website.String(), input.Label, input.CanRead, input.CanWrite, input.CanUpdate, input.CanUpgrade, input.CanDelete)
	if err != nil {
		return nil, err
	}

	existing, err := u.rbacRepo.FindRbacByLabelAndWebsite(input.Label, website.String())
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("rbac already exists for this label and website")
	}

	createdRbac, err := u.rbacRepo.CreateRbac(rbac)
	if err != nil {
		return nil, err
	}

	return createdRbac, nil
}

func (u *CreateRbacUseCase) GetByUUID(uuidStr string) (*domain.Rbac, error) {
	return u.rbacRepo.FindRbacByUUID(uuidStr)
}

func (u *CreateRbacUseCase) GetAll(websiteUUIDStr string) ([]*domain.Rbac, error) {
	return u.rbacRepo.GetRbacFromWebsite(websiteUUIDStr)
}
