package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/dtos"
	"github.com/ViitoJooj/verkoupe/internal/domain/usecases"
	"github.com/google/uuid"
)

type WebsiteController struct {
	createUseCase *usecases.CreateWebsiteUseCase
}

func NewWebsiteController(createUseCase *usecases.CreateWebsiteUseCase) *WebsiteController {
	return &WebsiteController{
		createUseCase: createUseCase,
	}
}

func (c *WebsiteController) Create(w http.ResponseWriter, r *http.Request) {
	websiteUUIDStr := r.Header.Get("X-Website-UUID")
	if websiteUUIDStr == "" {
		http.Error(w, "missing X-Website-UUID header", http.StatusBadRequest)
		return
	}

	_, err := uuid.Parse(websiteUUIDStr)
	if err != nil {
		http.Error(w, "invalid X-Website-UUID", http.StatusBadRequest)
		return
	}

	var req dtos.CreateWebsiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	website, err := c.createUseCase.Create(req.OwnerUUID, req.OwnerType, req.Label, req.URL, req.WriteIn, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := dtos.WebsiteResponse{
		UUID:        website.UUID.String(),
		OwnerUUID:   website.OwnerUUID.String(),
		OwnerType:   string(website.OwnerType),
		Label:       website.Label,
		URL:         website.URL,
		WriteIn:     string(website.WriteIn),
		Description: website.Description,
		CreatedAt:   website.CreatedAt.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
