package handlers

import (
	"encoding/json"
	"net/http"

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
	Id         string `json:"id"`
	Email      string `json:"email"`
	Updated_at string `json:"updated_at"`
	Created_at string `json:"created_at"`
}

type Response struct {
	Success  bool     `json:"success"`
	Message  string   `json:"message"`
	UserData UserData `json:"user"`
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

	resp := Response{
		Success: true,
		Message: "registered.",
		UserData: UserData{
			Id:         user.Id,
			Email:      user.Email,
			Created_at: user.Created_at,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := Response{
		Success: true,
		Message: "logged in.",
		UserData: UserData{
			Id:         user.Id,
			Email:      user.Email,
			Updated_at: user.Updated_at,
			Created_at: user.Created_at,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(resp)
}
