package usecases

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repositories/contracts"
	"github.com/google/uuid"
)

type CreateRbacUseCase struct {
	db       *sql.DB
	rbacRepo contracts.RbacContract
}

func NewCreateRbacUseCase(db *sql.DB, rbacRepo contracts.RbacContract) *CreateRbacUseCase {
	return &CreateRbacUseCase{
		db:       db,
		rbacRepo: rbacRepo,
	}
}

func (u *CreateRbacUseCase) Create(input *domain.Rbac, website uuid.UUID) (*domain.Rbac, error) {

	rbac, err := domain.NewRbac(website.String(), input.Label, input.CanRead, input.CanWrite, input.CanUpdate, input.CanUpgrade, input.CanDelete, u.db)
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
