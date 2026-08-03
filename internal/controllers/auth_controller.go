package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/dtos"
	"github.com/ViitoJooj/Jesterx/internal/usecases"
	"github.com/google/uuid"
)

type AuthController struct {
	registerUseCase *usecases.RegisterUserUseCase
}

func NewAuthController(registerUseCase *usecases.RegisterUserUseCase) *AuthController {
	return &AuthController{
		registerUseCase: registerUseCase,
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
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

	var req dtos.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := c.registerUseCase.Register(input, websiteUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := dtos.RegisterResponse{
		UUID:      user.UUID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
