package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/dtos"
	"github.com/ViitoJooj/Jesterx/internal/usecases"
)

type PlanController struct {
	createUseCase *usecases.CreatePlanUseCase
}

func NewPlanController(createUseCase *usecases.CreatePlanUseCase) *PlanController {
	return &PlanController{
		createUseCase: createUseCase,
	}
}

func (c *PlanController) Create(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	plan, err := c.createUseCase.Create(req.Name, req.Description, req.MaxWebsites, req.MaxRouters, req.MaxProducts, req.CostPerSaleRate, req.Coin, req.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := dtos.PlanResponse{
		UUID:            plan.UUID.String(),
		Name:            plan.Name,
		Description:     plan.Description,
		MaxWebsites:     plan.MaxWebsites,
		MaxRouters:      plan.MaxRouters,
		MaxProducts:     plan.MaxProducts,
		CostPerSaleRate: plan.CostPerSaleRate,
		Coin:            string(plan.Coin),
		Price:           plan.Price,
		CreatedAt:       plan.CreatedAt.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
