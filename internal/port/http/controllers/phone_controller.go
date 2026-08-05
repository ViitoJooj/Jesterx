package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/entities/enums"
	"github.com/ViitoJooj/verkoupe/internal/port/http/dtos"
	"github.com/ViitoJooj/verkoupe/internal/domain/repositories/contracts"
	"github.com/ViitoJooj/verkoupe/internal/domain/usecases"
	"github.com/google/uuid"
)

type PhoneController struct {
	createUseCase *usecases.CreatePhoneUseCase
	phoneRepo     contracts.PhoneContract
}

func NewPhoneController(createUseCase *usecases.CreatePhoneUseCase, phoneRepo contracts.PhoneContract) *PhoneController {
	return &PhoneController{
		createUseCase: createUseCase,
		phoneRepo:     phoneRepo,
	}
}

func (c *PhoneController) Create(w http.ResponseWriter, r *http.Request) {
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

	var req dtos.CreatePhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := &domain.Phone{
		OwnerUUID: uuid.MustParse(req.OwnerUUID),
		OwnerType: enums.OwnerType(req.OwnerType),
		Label:     req.Label,
		Number:    req.Number,
		IsDefault: req.IsDefault,
	}

	phone, err := c.createUseCase.Create(input, websiteUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := c.toResponse(phone)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (c *PhoneController) GetByUUID(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		http.Error(w, "missing uuid query parameter", http.StatusBadRequest)
		return
	}

	phone, err := c.phoneRepo.FindPhoneByUUID(uuidStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := c.toResponse(phone)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *PhoneController) GetAll(w http.ResponseWriter, r *http.Request) {
	websiteUUIDStr := r.Header.Get("X-Website-UUID")
	if websiteUUIDStr == "" {
		http.Error(w, "missing X-Website-UUID header", http.StatusBadRequest)
		return
	}

	phones, err := c.phoneRepo.GetPhonesFromWebsite(websiteUUIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]dtos.PhoneResponse, 0, len(phones))
	for _, phone := range phones {
		resp = append(resp, c.toResponse(phone))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *PhoneController) toResponse(phone *domain.Phone) dtos.PhoneResponse {
	updatedAt := ""
	if phone.UpdatedAt != nil {
		updatedAt = phone.UpdatedAt.String()
	}

	return dtos.PhoneResponse{
		UUID:        phone.UUID.String(),
		WebSiteUUID: phone.WebSiteUUID.String(),
		OwnerUUID:   phone.OwnerUUID.String(),
		OwnerType:   string(phone.OwnerType),
		Label:       phone.Label,
		Number:      phone.Number,
		IsDefault:   phone.IsDefault,
		UpdatedAt:   updatedAt,
		CreatedAt:   phone.CreatedAt.String(),
	}
}
