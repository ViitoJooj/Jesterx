package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/domain/enums"
	"github.com/ViitoJooj/Jesterx/internal/dtos"
	"github.com/ViitoJooj/Jesterx/internal/repositories/contracts"
	"github.com/ViitoJooj/Jesterx/internal/usecases"
	"github.com/google/uuid"
)

type TermsAcceptedController struct {
	createUseCase *usecases.CreateTermsAcceptedUseCase
	termsAccRepo  contracts.TermsAcceptedContract
}

func NewTermsAcceptedController(createUseCase *usecases.CreateTermsAcceptedUseCase, termsAccRepo contracts.TermsAcceptedContract) *TermsAcceptedController {
	return &TermsAcceptedController{
		createUseCase: createUseCase,
		termsAccRepo:  termsAccRepo,
	}
}

func (c *TermsAcceptedController) Create(w http.ResponseWriter, r *http.Request) {
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

	var req dtos.CreateTermsAcceptedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := &domain.TermsAcceptedBy{
		OwnerUUID: uuid.MustParse(req.OwnerUUID),
		OwnerType: enums.OwnerType(req.OwnerType),
	}

	termsAccepted, err := c.createUseCase.Create(input, websiteUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := c.toResponse(termsAccepted)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (c *TermsAcceptedController) GetByUUID(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		http.Error(w, "missing uuid query parameter", http.StatusBadRequest)
		return
	}

	termsAccepted, err := c.termsAccRepo.FindTermsAcceptedByUUID(uuidStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := c.toResponse(termsAccepted)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *TermsAcceptedController) GetAll(w http.ResponseWriter, r *http.Request) {
	websiteUUIDStr := r.Header.Get("X-Website-UUID")
	if websiteUUIDStr == "" {
		http.Error(w, "missing X-Website-UUID header", http.StatusBadRequest)
		return
	}

	termsAccepteds, err := c.termsAccRepo.GetTermsAcceptedFromWebsite(websiteUUIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]dtos.TermsAcceptedResponse, 0, len(termsAccepteds))
	for _, ta := range termsAccepteds {
		resp = append(resp, c.toResponse(ta))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *TermsAcceptedController) toResponse(ta *domain.TermsAcceptedBy) dtos.TermsAcceptedResponse {
	return dtos.TermsAcceptedResponse{
		UUID:         ta.UUID.String(),
		WebSiteUUID:  ta.WebSiteUUID.String(),
		OwnerUUID:    ta.OwnerUUID.String(),
		OwnerType:    string(ta.OwnerType),
		AcceptedWhen: ta.AcceptedWhen.String(),
	}
}
