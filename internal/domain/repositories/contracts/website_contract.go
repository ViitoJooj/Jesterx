package contracts

import (
	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
)

type WebsiteContract interface {
	CreateWebsite(website *domain.Website) (*domain.Website, error)
	FindWebsiteByUUID(uuid string) (*domain.Website, error)
	FindWebsiteByLabel(label string) (*domain.Website, error)
	FindWebsitesByOwner(ownerUUID string) ([]*domain.Website, error)
	GetWebsites() ([]*domain.Website, error)
	UpdateWebsiteByUUID(uuid string) error
	DeleteWebsiteByUUID(uuid string) error
	DeleteWebsitesByUUIDS(uuid []string) error
}
