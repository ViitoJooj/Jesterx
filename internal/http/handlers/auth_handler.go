package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/security"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

type RegisterRequest struct {
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserData struct {
	Id         string    `json:"id"`
	Email      string    `json:"email"`
	Updated_at time.Time `json:"updated_at"`
	Created_at time.Time `json:"created_at"`
}

type AuthResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}

type ResponseRefreshToken struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type UserMeResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(
		req.First_name,
		req.Last_name,
		req.Email,
		req.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims := security.RefreshTokenClaims{
		Iss: "https://jesterx.com.br",
		Sub: user.Id,
		Exp: time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	refreshToken, err := security.RefreshToken(claims)
	if err != nil {
		http.Error(w, "internal error", http.StatusBadGateway)
		return
	}

	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   2592000,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	resp := AuthResponse{
		Success: true,
		Message: "registered.",
		Data: UserData{
			Id:         user.Id,
			Email:      user.Email,
			Created_at: user.Created_at,
		},
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims := security.RefreshTokenClaims{
		Iss: "https://jesterx.com.br",
		Sub: user.Id,
		Exp: time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	refreshToken, err := security.RefreshToken(claims)
	if err != nil {
		http.Error(w, "Internal error", http.StatusBadGateway)
		return
	}

	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   2592000,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	resp := AuthResponse{
		Success: true,
		Message: "logged in.",
		Data: UserData{
			Id:         user.Id,
			Email:      user.Email,
			Updated_at: user.Updated_at,
			Created_at: user.Created_at,
		},
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	accessToken, err := h.authService.Refresh(refreshCookie.Value)
	if err != nil {
		http.Error(w, "not allowed", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   900,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	resp := ResponseRefreshToken{
		Success: true,
		Message: "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	accessCookie, err := r.Cookie("access_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	user, err := h.authService.Me(accessCookie.Value)
	if err != nil {
		http.Error(w, "not allowed", http.StatusUnauthorized)
		return
	}

	resp := UserMeResponse{
		ID:        user.Id,
		FirstName: user.First_name,
		LastName:  user.Last_name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.Created_at.Format(time.RFC3339),
		UpdatedAt: user.Updated_at.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
