package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/port/http/dtos"
	"github.com/ViitoJooj/verkoupe/internal/domain/usecases"
	"github.com/google/uuid"
)

type OrganizationController struct {
	createUseCase *usecases.CreateOrganizationUseCase
}

func NewOrganizationController(createUseCase *usecases.CreateOrganizationUseCase) *OrganizationController {
	return &OrganizationController{
		createUseCase: createUseCase,
	}
}

func (c *OrganizationController) Create(w http.ResponseWriter, r *http.Request) {
	websiteUUIDStr := r.Header.Get("X-Website-UUID")
	if websiteUUIDStr == "" {
		http.Error(w, "missing X-Website-UUID header", http.StatusBadRequest)
		return
	}

	websiteUUID, err := uuid.Parse(websiteUUIDStr)
	if err != nil {
		http.Error(w, "invalid X-Website-UUID", http.StatusBadRequest)
		return
	}

	var req dtos.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := &domain.OrganizationBR{
		OwnerUUID: uuid.MustParse(req.OwnerUUID),
		ImageURL:  req.ImageURL,
		Name:      req.Name,
		TradeName: req.TradeName,
		CNPJ:      req.CNPJ,
	}

	org, err := c.createUseCase.Create(input, websiteUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := c.toResponse(org)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (c *OrganizationController) GetByUUID(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		http.Error(w, "missing uuid query parameter", http.StatusBadRequest)
		return
	}

	org, err := c.createUseCase.GetByUUID(uuidStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := c.toResponse(org)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *OrganizationController) GetAll(w http.ResponseWriter, r *http.Request) {
	websiteUUIDStr := r.Header.Get("X-Website-UUID")
	if websiteUUIDStr == "" {
		http.Error(w, "missing X-Website-UUID header", http.StatusBadRequest)
		return
	}

	orgs, err := c.createUseCase.GetAll(websiteUUIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]dtos.OrganizationResponse, 0, len(orgs))
	for _, org := range orgs {
		resp = append(resp, c.toResponse(org))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *OrganizationController) toResponse(org *domain.OrganizationBR) dtos.OrganizationResponse {
	updatedAt := ""
	if org.UpdatedAt != nil {
		updatedAt = org.UpdatedAt.String()
	}

	return dtos.OrganizationResponse{
		UUID:        org.UUID.String(),
		WebSiteUUID: org.WebSiteUUID.String(),
		OwnerUUID:   org.OwnerUUID.String(),
		ImageURL:    org.ImageURL,
		Name:        org.Name,
		TradeName:   org.TradeName,
		CNPJ:        org.CNPJ,
		UpdatedAt:   updatedAt,
		CreatedAt:   org.CreatedAt.String(),
	}
}
