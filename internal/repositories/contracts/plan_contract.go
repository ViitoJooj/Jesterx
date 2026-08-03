package contracts

import (
	"github.com/ViitoJooj/Jesterx/internal/domain"
)

type PlanContract interface {
	CreatePlan(plan *domain.JesterxPlans) (*domain.JesterxPlans, error)
	FindPlanByUUID(uuid string) (*domain.JesterxPlans, error)
	FindPlanByName(name string) (*domain.JesterxPlans, error)
	GetPlans() ([]*domain.JesterxPlans, error)
	UpdatePlanByUUID(uuid string) error
	DeletePlanByUUID(uuid string) error
	DeletePlansByUUIDS(uuid []string) error
}
