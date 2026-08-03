package contracts

import (
	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
)

type OrganizationContract interface {
	CreateOrganization(org *domain.OrganizationBR) (*domain.OrganizationBR, error)
	FindOrganizationByUUID(uuid string) (*domain.OrganizationBR, error)
	FindOrganizationByCNPJAndWebsite(cnpj string, websiteUUID string) (*domain.OrganizationBR, error)
	GetOrganizationsFromWebsite(websiteUUID string) ([]*domain.OrganizationBR, error)
	GetOrganizations() ([]*domain.OrganizationBR, error)
	UpdateOrganizationByUUID(uuid string) error
	DeleteOrganizationByUUID(uuid string) error
	DeleteOrganizationsByUUIDS(uuid []string) error
}
