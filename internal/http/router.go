package http

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/http/handlers"
)

func NewRouter() *http.ServeMux {
	return http.NewServeMux()
}

func RegisterAuthRoutes(mux *http.ServeMux, h *handlers.AuthHandler) {
	mux.HandleFunc("POST /api/v1/auth/register", h.Register)
	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
	mux.HandleFunc("GET /api/v1/auth/refresh", h.Refresh)
	mux.HandleFunc("GET /api/v1/auth/me", h.Me)
}

func RegisterWebsiteRoutes(mux *http.ServeMux, h *handlers.WebSiteHandler) {
	mux.HandleFunc("POST /api/v1/websites", h.CreateWebSite)
}
