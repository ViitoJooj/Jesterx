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

type AddressController struct {
	createUseCase *usecases.CreateAddressUseCase
	addressRepo   contracts.AddressContract
}

func NewAddressController(createUseCase *usecases.CreateAddressUseCase, addressRepo contracts.AddressContract) *AddressController {
	return &AddressController{
		createUseCase: createUseCase,
		addressRepo:   addressRepo,
	}
}

func (c *AddressController) Create(w http.ResponseWriter, r *http.Request) {
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

	var req dtos.CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := &domain.AddressBR{
		OwnerUUID:      uuid.MustParse(req.OwnerUUID),
		OwnerType:      enums.OwnerType(req.OwnerType),
		Label:          req.Label,
		AddressLine1:   req.AddressLine1,
		AddressLine2:   req.AddressLine2,
		Neighborhood:   req.Neighborhood,
		City:           req.City,
		State:          req.State,
		StateCode:      req.StateCode,
		PostalCode:     req.PostalCode,
		ReferencePoint: req.ReferencePoint,
		DeliveryNotes:  req.DeliveryNotes,
		IsDefault:      req.IsDefault,
	}

	address, err := c.createUseCase.Create(input, websiteUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := c.toResponse(address)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (c *AddressController) GetByUUID(w http.ResponseWriter, r *http.Request) {
	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		http.Error(w, "missing uuid query parameter", http.StatusBadRequest)
		return
	}

	address, err := c.addressRepo.FindAddressByUUID(uuidStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := c.toResponse(address)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *AddressController) GetAll(w http.ResponseWriter, r *http.Request) {
	websiteUUIDStr := r.Header.Get("X-Website-UUID")
	if websiteUUIDStr == "" {
		http.Error(w, "missing X-Website-UUID header", http.StatusBadRequest)
		return
	}

	addresses, err := c.addressRepo.GetAddressesFromWebsite(websiteUUIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]dtos.AddressResponse, 0, len(addresses))
	for _, address := range addresses {
		resp = append(resp, c.toResponse(address))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (c *AddressController) toResponse(address *domain.AddressBR) dtos.AddressResponse {
	updatedAt := ""
	if address.UpdatedAt != nil {
		updatedAt = address.UpdatedAt.String()
	}

	return dtos.AddressResponse{
		UUID:           address.UUID.String(),
		WebSiteUUID:    address.WebSiteUUID.String(),
		OwnerUUID:      address.OwnerUUID.String(),
		OwnerType:      string(address.OwnerType),
		Label:          address.Label,
		AddressLine1:   address.AddressLine1,
		AddressLine2:   address.AddressLine2,
		Neighborhood:   address.Neighborhood,
		City:           address.City,
		State:          address.State,
		StateCode:      address.StateCode,
		PostalCode:     address.PostalCode,
		ReferencePoint: address.ReferencePoint,
		DeliveryNotes:  address.DeliveryNotes,
		IsDefault:      address.IsDefault,
		UpdatedAt:      updatedAt,
		CreatedAt:      address.CreatedAt.String(),
	}
}
