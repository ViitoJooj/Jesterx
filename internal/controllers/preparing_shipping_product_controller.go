package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/dtos"
	"github.com/ViitoJooj/Jesterx/internal/usecases"
)

type PreparingShippingProductController struct {
	createUseCase *usecases.CreatePreparingShippingProductUseCase
}

func NewPreparingShippingProductController(createUseCase *usecases.CreatePreparingShippingProductUseCase) *PreparingShippingProductController {
	return &PreparingShippingProductController{
		createUseCase: createUseCase,
	}
}

func (c *PreparingShippingProductController) Create(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreatePreparingShippingProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	psp, err := c.createUseCase.Create(req.ProductUUID, req.AddressUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := dtos.PreparingShippingProductResponse{
		UUID:        psp.UUID.String(),
		ProductUUID: psp.ProductUUID.String(),
		AddressUUID: psp.AddressUUID.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
