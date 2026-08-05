package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/dtos"
	"github.com/ViitoJooj/verkoupe/internal/domain/usecases"
)

type StorageProductController struct {
	createUseCase *usecases.CreateStorageProductUseCase
}

func NewStorageProductController(createUseCase *usecases.CreateStorageProductUseCase) *StorageProductController {
	return &StorageProductController{
		createUseCase: createUseCase,
	}
}

func (c *StorageProductController) Create(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateStorageProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sp, err := c.createUseCase.Create(req.ProductUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := dtos.StorageProductResponse{
		UUID:        sp.UUID.String(),
		ProductUUID: sp.ProductUUID.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
