package contracts

import (
	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
)

type WebsiteComponentContract interface {
	CreateWebsiteComponent(component *domain.ComponentWebsites) (*domain.ComponentWebsites, error)
	FindWebsiteComponentByUUID(uuid string) (*domain.ComponentWebsites, error)
	FindWebsiteComponentByPath(path string) (*domain.ComponentWebsites, error)
	GetWebsiteComponentsFromWebsite(websiteUUID string) ([]*domain.ComponentWebsites, error)
	GetWebsiteComponents() ([]*domain.ComponentWebsites, error)
	UpdateWebsiteComponentByUUID(uuid string) error
	DeleteWebsiteComponentByUUID(uuid string) error
	DeleteWebsiteComponentsByUUIDS(uuid []string) error
}
