package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/port/http/dtos"
	"github.com/ViitoJooj/Jesterx/internal/domain/usecases"
)

type ProductTagController struct {
	createUseCase *usecases.CreateProductTagUseCase
}

func NewProductTagController(createUseCase *usecases.CreateProductTagUseCase) *ProductTagController {
	return &ProductTagController{
		createUseCase: createUseCase,
	}
}

func (c *ProductTagController) Create(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateProductTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tag, err := c.createUseCase.Create(req.ProductUUID, req.Label)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := dtos.ProductTagResponse{
		UUID:        tag.UUID.String(),
		ProductUUID: tag.ProductUUID.String(),
		Label:       tag.Label,
		CreatedAt:   tag.CreatedAt.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
