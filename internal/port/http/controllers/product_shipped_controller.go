package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/port/http/dtos"
	"github.com/ViitoJooj/Jesterx/internal/domain/usecases"
)

type ProductShippedController struct {
	createUseCase *usecases.CreateProductShippedUseCase
}

func NewProductShippedController(createUseCase *usecases.CreateProductShippedUseCase) *ProductShippedController {
	return &ProductShippedController{
		createUseCase: createUseCase,
	}
}

func (c *ProductShippedController) Create(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateProductShippedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	productShipped, err := c.createUseCase.Create(req.ProductUUID, req.AddressUUID, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := dtos.ProductShippedResponse{
		UUID:        productShipped.UUID.String(),
		ProductUUID: productShipped.ProductUUID.String(),
		AddressUUID: productShipped.AddressUUID.String(),
		Status:      productShipped.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
