package usecases

import (
	"database/sql"
	"errors"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/google/uuid"
)

type CreateOrganizationUseCase struct {
	db      *sql.DB
	orgRepo contracts.OrganizationContract
}

func NewCreateOrganizationUseCase(db *sql.DB, orgRepo contracts.OrganizationContract) *CreateOrganizationUseCase {
	return &CreateOrganizationUseCase{
		db:      db,
		orgRepo: orgRepo,
	}
}

func (u *CreateOrganizationUseCase) Create(input *domain.OrganizationBR, website uuid.UUID) (*domain.OrganizationBR, error) {

	org, err := domain.NewOrganization(website.String(), input.OwnerUUID.String(), input.ImageURL, input.Name, input.TradeName, input.CNPJ, u.db)
	if err != nil {
		return nil, err
	}

	existing, err := u.orgRepo.FindOrganizationByCNPJAndWebsite(input.CNPJ, website.String())
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("organization already exists for this CNPJ and website")
	}

	createdOrg, err := u.orgRepo.CreateOrganization(org)
	if err != nil {
		return nil, err
	}

	return createdOrg, nil
}
