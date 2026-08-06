package controllers

import (
	"encoding/json"
	"net/http"

	domain "github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/usecases"
	"github.com/ViitoJooj/verkoupe/internal/port/http/dtos"
	"github.com/ViitoJooj/verkoupe/internal/port/http/middleware"
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
	var req dtos.CreateWebsiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse("RAX-004", "invalid request body"))
		return
	}

	website, err := c.createUseCase.Create(req.OwnerUUID, req.OwnerType, req.Label, req.URL, req.WriteIn, req.Description)
	if err != nil {
		writeJSON(w, http.StatusConflict, errorResponse("RDX-002", err.Error()))
		return
	}

	writeJSON(w, http.StatusCreated, websiteToResponse(website))
}

func (c *WebsiteController) GetByUUID(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.PathValue("uuid")
	if uuidStr == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse("RDX-003", "missing uuid"))
		return
	}

	website, err := c.createUseCase.GetByUUID(uuidStr)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse("RDX-001", "website not found"))
		return
	}

	writeJSON(w, http.StatusOK, websiteToResponse(website))
}

func (c *WebsiteController) ListByOwner(w http.ResponseWriter, r *http.Request) {
	userUUID := middleware.GetUserUUID(r)
	if userUUID == "" {
		writeJSON(w, http.StatusUnauthorized, errorResponse("RBX-012", "unauthorized"))
		return
	}

	websites, err := c.createUseCase.ListByOwner(userUUID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse("RAX-001", "internal error"))
		return
	}

	if websites == nil {
		websites = []*domain.Website{}
	}

	responses := make([]dtos.WebsiteResponse, 0, len(websites))
	for _, w := range websites {
		responses = append(responses, websiteToResponse(w))
	}

	writeJSON(w, http.StatusOK, responses)
}

func (c *WebsiteController) ListAll(w http.ResponseWriter, r *http.Request) {
	websites, err := c.createUseCase.ListAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse("RAX-001", "internal error"))
		return
	}

	if websites == nil {
		websites = []*domain.Website{}
	}

	responses := make([]dtos.WebsiteResponse, 0, len(websites))
	for _, w := range websites {
		responses = append(responses, websiteToResponse(w))
	}

	writeJSON(w, http.StatusOK, responses)
}

func (c *WebsiteController) Update(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.PathValue("uuid")
	if uuidStr == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse("RDX-003", "missing uuid"))
		return
	}

	if _, err := uuid.Parse(uuidStr); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse("RDX-003", "invalid uuid"))
		return
	}

	if _, err := c.createUseCase.GetByUUID(uuidStr); err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse("RDX-001", "website not found"))
		return
	}

	if err := c.createUseCase.Update(uuidStr); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse("RAX-001", "internal error"))
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (c *WebsiteController) Delete(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.PathValue("uuid")
	if uuidStr == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse("RDX-003", "missing uuid"))
		return
	}

	if err := c.createUseCase.Delete(uuidStr); err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse("RDX-001", "website not found"))
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func websiteToResponse(w *domain.Website) dtos.WebsiteResponse {
	return dtos.WebsiteResponse{
		UUID:        w.UUID.String(),
		OwnerUUID:   w.OwnerUUID.String(),
		OwnerType:   string(w.OwnerType),
		Label:       w.Label,
		URL:         w.URL,
		WriteIn:     string(w.WriteIn),
		Description: w.Description,
		CreatedAt:   w.CreatedAt.String(),
	}
}
