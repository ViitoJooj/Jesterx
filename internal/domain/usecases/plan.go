package usecases

import (

	domain "github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
)

type CreatePlanUseCase struct {
	repository contracts.PlanContract
}

func NewCreatePlanUseCase(repository contracts.PlanContract) *CreatePlanUseCase {
	return &CreatePlanUseCase{repository: repository}
}

func (u *CreatePlanUseCase) Create(name string, description string, maxWebsites int, maxRouters int, maxProducts int, costPerSaleRate int, coin string, price int) (*domain.VerkoupePlan, error) {
	plan, err := domain.NewPlan(name, description, maxWebsites, maxRouters, maxProducts, costPerSaleRate, coin, price)
	if err != nil {
		return nil, err
	}
	return u.repository.CreatePlan(plan)
}
