package contracts

import (
	"github.com/ViitoJooj/Jesterx/internal/domain/entities"
)

type RbacContract interface {
	CreateRbac(rbac *domain.Rbac) (*domain.Rbac, error)
	FindRbacByUUID(uuid string) (*domain.Rbac, error)
	FindRbacByLabelAndWebsite(label string, websiteUUID string) (*domain.Rbac, error)
	GetRbacFromWebsite(websiteUUID string) ([]*domain.Rbac, error)
	GetRbac() ([]*domain.Rbac, error)
	UpdateRbacByUUID(uuid string) error
	DeleteRbacByUUID(uuid string) error
	DeleteRbacByUUIDS(uuid []string) error
}
