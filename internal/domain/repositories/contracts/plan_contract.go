package contracts

import (
	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
)

type PlanContract interface {
	CreatePlan(plan *domain.verkoupePlans) (*domain.verkoupePlans, error)
	FindPlanByUUID(uuid string) (*domain.verkoupePlans, error)
	FindPlanByName(name string) (*domain.verkoupePlans, error)
	GetPlans() ([]*domain.verkoupePlans, error)
	UpdatePlanByUUID(uuid string) error
	DeletePlanByUUID(uuid string) error
	DeletePlansByUUIDS(uuid []string) error
}
