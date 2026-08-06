package usecases

import (
	"errors"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/google/uuid"
)

type CreateOrganizationUseCase struct {
	orgRepo contracts.OrganizationContract
}

func NewCreateOrganizationUseCase(orgRepo contracts.OrganizationContract) *CreateOrganizationUseCase {
	return &CreateOrganizationUseCase{
		orgRepo: orgRepo,
	}
}

func (u *CreateOrganizationUseCase) Create(input *domain.OrganizationBR, website uuid.UUID) (*domain.OrganizationBR, error) {

	org, err := domain.NewOrganization(website.String(), input.OwnerUUID.String(), input.ImageURL, input.Name, input.TradeName, input.CNPJ)
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

func (u *CreateOrganizationUseCase) GetByUUID(uuidStr string) (*domain.OrganizationBR, error) {
	return u.orgRepo.FindOrganizationByUUID(uuidStr)
}

func (u *CreateOrganizationUseCase) GetAll(websiteUUIDStr string) ([]*domain.OrganizationBR, error) {
	return u.orgRepo.GetOrganizationsFromWebsite(websiteUUIDStr)
}
