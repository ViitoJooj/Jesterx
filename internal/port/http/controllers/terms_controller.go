package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/port/http/dtos"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/ViitoJooj/verkoupe/internal/domain/usecases"
)

type TermsController struct {
	createUseCase *usecases.CreateTermsUseCase
	termsRepo     contracts.TermsContract
}

func NewTermsController(createUseCase *usecases.CreateTermsUseCase, termsRepo contracts.TermsContract) *TermsController {
	return &TermsController{
		createUseCase: createUseCase,
		termsRepo:     termsRepo,
	}
}

func (c *TermsController) Create(w http.ResponseWriter, r *http.Request) {
	websiteUUIDStr := r.Header.Get("X-Website-UUID")
	if websiteUUIDStr == "" {
		http.Error(w, "missing X-Website-UUID header", http.StatusBadRequest)
		return
	}

	var req dtos.CreateTermsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := &domain.Terms{
		Name:        req.Name,
		Description: req.Description,
	}

	terms, err := c.createUseCase.Create(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := c.toResponse(terms)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

	_ = websiteUUIDStr
}

func (c *TermsController) GetByUUID(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		http.Error(w, "missing uuid query parameter", http.StatusBadRequest)
		return
	}

	terms, err := c.termsRepo.FindTermsByUUID(uuidStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := c.toResponse(terms)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *TermsController) GetAll(w http.ResponseWriter, r *http.Request) {
	terms, err := c.termsRepo.GetTerms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]dtos.TermsResponse, 0, len(terms))
	for _, t := range terms {
		resp = append(resp, c.toResponse(t))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *TermsController) toResponse(terms *domain.Terms) dtos.TermsResponse {
	updatedAt := ""
	if terms.UpdatedAt != nil {
		updatedAt = terms.UpdatedAt.String()
	}

	return dtos.TermsResponse{
		UUID:        terms.UUID.String(),
		Name:        terms.Name,
		Description: terms.Description,
		UpdatedAt:   updatedAt,
		CreatedAt:   terms.CreatedAt.String(),
	}
}
