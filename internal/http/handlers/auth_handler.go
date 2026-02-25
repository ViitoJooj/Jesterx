package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/config"
	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
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
	WebsiteId  string    `json:"website_id"`
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

	websiteId := r.Header.Get("X-Website-Id")
	if websiteId == "" {
		http.Error(w, "invalid website id", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(
		websiteId,
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
		Iss:       "https://jesterx.com.br",
		Sub:       user.Id,
		WebsiteId: user.WebsiteId,
		Exp:       time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	refreshToken, err := security.RefreshToken(claims)
	if err != nil {
		http.Error(w, "internal error", http.StatusBadGateway)
		return
	}

	cookie := http.Cookie{
		Name:     security.RefreshCookieName(user.WebsiteId),
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		Secure:   !config.IsDev,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	resp := AuthResponse{
		Success: true,
		Message: "registered.",
		Data: UserData{
			Id:         user.Id,
			WebsiteId:  user.WebsiteId,
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

	websiteId := r.Header.Get("X-Website-Id")
	if websiteId == "" {
		http.Error(w, "invalid website id", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(
		websiteId,
		req.Email,
		req.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	claims := security.RefreshTokenClaims{
		Iss:       "https://jesterx.com.br",
		Sub:       user.Id,
		WebsiteId: user.WebsiteId,
		Exp:       time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	refreshToken, err := security.RefreshToken(claims)
	if err != nil {
		http.Error(w, "Internal error", http.StatusBadGateway)
		return
	}

	cookie := http.Cookie{
		Name:     security.RefreshCookieName(user.WebsiteId),
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		Secure:   !config.IsDev,
		SameSite: http.SameSiteLaxMode,
	}

	resp := AuthResponse{
		Success: true,
		Message: "logged in.",
		Data: UserData{
			Id:         user.Id,
			WebsiteId:  user.WebsiteId,
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

	websiteId := r.Header.Get("X-Website-Id")
	if websiteId == "" {
		http.Error(w, "invalid website id", http.StatusBadRequest)
		return
	}

	refreshCookie, err := r.Cookie(security.RefreshCookieName(websiteId))
	if err != nil {
		http.Error(w, "refresh token missing", http.StatusUnauthorized)
		return
	}

	accessToken, err := h.authService.Refresh(refreshCookie.Value)
	if err != nil {
		log.Println(err)
		http.Error(w, "not allowed", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     security.AccessCookieName(websiteId),
		Value:    accessToken,
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
		Secure:   !config.IsDev,
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
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.Me(userID)
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

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	websiteId := r.Header.Get("X-Website-Id")
	if websiteId == "" {
		http.Error(w, "invalid website id", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   security.RefreshCookieName(websiteId),
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:   security.AccessCookieName(websiteId),
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusNoContent)
}
