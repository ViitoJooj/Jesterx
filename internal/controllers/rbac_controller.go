package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/dtos"
	"github.com/ViitoJooj/Jesterx/internal/repositories/contracts"
	"github.com/ViitoJooj/Jesterx/internal/usecases"
	"github.com/google/uuid"
)

type RbacController struct {
	createUseCase *usecases.CreateRbacUseCase
	rbacRepo      contracts.RbacContract
}

func NewRbacController(createUseCase *usecases.CreateRbacUseCase, rbacRepo contracts.RbacContract) *RbacController {
	return &RbacController{
		createUseCase: createUseCase,
		rbacRepo:      rbacRepo,
	}
}

func (c *RbacController) Create(w http.ResponseWriter, r *http.Request) {
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

	var req dtos.CreateRbacRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := &domain.Rbac{
		Label:      req.Label,
		CanRead:    req.CanRead,
		CanWrite:   req.CanWrite,
		CanUpdate:  req.CanUpdate,
		CanUpgrade: req.CanUpgrade,
		CanDelete:  req.CanDelete,
	}

	rbac, err := c.createUseCase.Create(input, websiteUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := c.toResponse(rbac)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (c *RbacController) GetByUUID(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		http.Error(w, "missing uuid query parameter", http.StatusBadRequest)
		return
	}

	rbac, err := c.rbacRepo.FindRbacByUUID(uuidStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := c.toResponse(rbac)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *RbacController) GetAll(w http.ResponseWriter, r *http.Request) {
	websiteUUIDStr := r.Header.Get("X-Website-UUID")
	if websiteUUIDStr == "" {
		http.Error(w, "missing X-Website-UUID header", http.StatusBadRequest)
		return
	}

	rbacs, err := c.rbacRepo.GetRbacFromWebsite(websiteUUIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]dtos.RbacResponse, 0, len(rbacs))
	for _, rbac := range rbacs {
		resp = append(resp, c.toResponse(rbac))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *RbacController) toResponse(rbac *domain.Rbac) dtos.RbacResponse {
	updatedAt := ""
	if rbac.UpdatedAt != nil {
		updatedAt = rbac.UpdatedAt.String()
	}

	return dtos.RbacResponse{
		UUID:        rbac.UUID.String(),
		WebSiteUUID: rbac.WebSiteUUID.String(),
		Label:       rbac.Label,
		CanRead:     rbac.CanRead,
		CanWrite:    rbac.CanWrite,
		CanUpdate:   rbac.CanUpdate,
		CanUpgrade:  rbac.CanUpgrade,
		CanDelete:   rbac.CanDelete,
		UpdatedAt:   updatedAt,
		CreatedAt:   rbac.CreatedAt.String(),
	}
}
