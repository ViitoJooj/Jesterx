package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/port/http/dtos"
	"github.com/ViitoJooj/Jesterx/internal/domain/usecases"
)

type CupomController struct {
	createUseCase *usecases.CreateCupomUseCase
}

func NewCupomController(createUseCase *usecases.CreateCupomUseCase) *CupomController {
	return &CupomController{
		createUseCase: createUseCase,
	}
}

func (c *CupomController) Create(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateCupomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cupom, err := c.createUseCase.Create(req.TagUUID, req.Label, req.Description, req.Value, req.ValueType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := dtos.CupomResponse{
		UUID:        cupom.UUID.String(),
		TagUUID:     cupom.TagUUID.String(),
		Label:       cupom.Label,
		Description: cupom.Description,
		Value:       cupom.Value,
		ValueType:   string(cupom.ValueType),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
