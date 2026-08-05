package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/dtos"
	"github.com/ViitoJooj/verkoupe/internal/domain/usecases"
)

type ProductController struct {
	createUseCase *usecases.CreateProductUseCase
}

func NewProductController(createUseCase *usecases.CreateProductUseCase) *ProductController {
	return &ProductController{
		createUseCase: createUseCase,
	}
}

func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	product, err := c.createUseCase.Create(req.Name, req.Description, req.ShortDescription, req.Height, req.Width, req.Thickness, req.Active)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := dtos.ProductResponse{
		UUID:             product.UUID.String(),
		Name:             product.Name,
		Description:      product.Description,
		ShortDescription: product.ShortDescription,
		Active:           product.Active,
		CreatedAt:        product.CreatedAt.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
