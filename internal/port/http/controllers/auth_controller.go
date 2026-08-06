package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	domain "github.com/ViitoJooj/verkoupe/internal/domain/entities"
	"github.com/ViitoJooj/verkoupe/internal/domain/usecases"
	"github.com/ViitoJooj/verkoupe/internal/port/http/dtos"
	"github.com/ViitoJooj/verkoupe/pkg/token"
	"github.com/google/uuid"
)

type AuthController struct {
	authUseCase *usecases.AuthUseCase
}

func NewAuthController(authUseCase *usecases.AuthUseCase) *AuthController {
	return &AuthController{
		authUseCase: authUseCase,
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	websiteUUIDStr := r.Header.Get("X-Website-UUID")
	if websiteUUIDStr == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse("RAX-003", "missing X-Website-UUID header"))
		return
	}

	websiteUUID, err := uuid.Parse(websiteUUIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse("RAX-003", "invalid X-Website-UUID"))
		return
	}

	var req dtos.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse("RAX-004", "invalid request body"))
		return
	}

	input := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := c.authUseCase.Register(input, websiteUUID)
	if err != nil {
		writeJSON(w, http.StatusConflict, errorResponse("RBX-005", err.Error()))
		return
	}

	resp := dtos.RegisterResponse{
		UUID:      user.UUID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req dtos.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse("RAX-004", "invalid request body"))
		return
	}

	websiteUUIDStr := r.Header.Get("X-Website-UUID")
	if websiteUUIDStr == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse("RAX-003", "missing X-Website-UUID header"))
		return
	}

	user, err := c.authUseCase.Login(req.Email, req.Password, websiteUUIDStr)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse("RBX-001", "invalid credentials"))
		return
	}

	pasetoSecret := os.Getenv("PASETO_SECRET_KEY")
	if pasetoSecret == "" {
		pasetoSecret = "dev-secret-key-change-in-production-32b"
	}

	accessToken, err := token.GenerateAccess(
		[]byte(pasetoSecret),
		user.UUID.String(),
		user.WebSiteUUID.String(),
		user.Role,
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse("RAX-001", "could not generate token"))
		return
	}

	resp := dtos.LoginResponse{
		Token:       accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   900,
		UserUUID:    user.UUID.String(),
		WebsiteUUID: user.WebSiteUUID.String(),
		Name:        user.Name,
		Email:       user.Email,
	}

	writeJSON(w, http.StatusOK, resp)
}
