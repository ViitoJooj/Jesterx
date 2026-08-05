package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/dtos"
	"github.com/ViitoJooj/verkoupe/internal/domain/usecases"
	"github.com/google/uuid"
)

type WebsiteComponentController struct {
	createUseCase *usecases.CreateWebsiteComponentUseCase
}

func NewWebsiteComponentController(createUseCase *usecases.CreateWebsiteComponentUseCase) *WebsiteComponentController {
	return &WebsiteComponentController{
		createUseCase: createUseCase,
	}
}

func (c *WebsiteComponentController) Create(w http.ResponseWriter, r *http.Request) {
	websiteUUIDStr := r.Header.Get("X-Website-UUID")
	if websiteUUIDStr == "" {
		http.Error(w, "missing X-Website-UUID header", http.StatusBadRequest)
		return
	}

	tenantWebsiteUUID, err := uuid.Parse(websiteUUIDStr)
	if err != nil {
		http.Error(w, "invalid X-Website-UUID", http.StatusBadRequest)
		return
	}

	var req dtos.CreateWebsiteComponentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	component, err := c.createUseCase.Create(req.WebsiteUUID, req.LogoURL, req.Tittle, req.Description, req.Path, req.Content, req.Visits, tenantWebsiteUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := dtos.WebsiteComponentResponse{
		UUID:        component.UUID.String(),
		WebsiteUUID: component.WebsiteUUID.String(),
		LogoURL:     component.LogoURL,
		Tittle:      component.Tittle,
		Description: component.Description,
		Path:        component.Path,
		Content:     component.Content,
		Visits:      component.Visists,
		CreatedAt:   component.CreatedAt.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
