package usecases

import (
	"database/sql"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repositories/contracts"
)

type CreatePlanUseCase struct {
	db   *sql.DB
	repo contracts.PlanContract
}

func NewCreatePlanUseCase(db *sql.DB, repo contracts.PlanContract) *CreatePlanUseCase {
	return &CreatePlanUseCase{db: db, repo: repo}
}

func (u *CreatePlanUseCase) Create(name string, description string, maxWebsites int, maxRouters int, maxProducts int, costPerSaleRate int, coin string, price int) (*domain.JesterxPlans, error) {
	plan, err := domain.NewPlan(name, description, maxWebsites, maxRouters, maxProducts, costPerSaleRate, coin, price)
	if err != nil {
		return nil, err
	}
	return u.repo.CreatePlan(plan)
}
