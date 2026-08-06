package contracts

import (
	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
)

type PlanContract interface {
	CreatePlan(plan *domain.VerkoupePlan) (*domain.VerkoupePlan, error)
	FindPlanByUUID(uuid string) (*domain.VerkoupePlan, error)
	FindPlanByName(name string) (*domain.VerkoupePlan, error)
	GetPlans() ([]*domain.VerkoupePlan, error)
	UpdatePlanByUUID(uuid string) error
	DeletePlanByUUID(uuid string) error
	DeletePlansByUUIDS(uuid []string) error
}
