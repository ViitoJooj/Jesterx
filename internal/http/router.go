package http

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/http/handlers"
	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

func NewRouter() *http.ServeMux {
	return http.NewServeMux()
}

func RegisterAuthRoutes(mux *http.ServeMux, h *handlers.AuthHandler, authService *service.AuthService) {
	mux.HandleFunc("POST /api/v1/auth/register", h.Register)
	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
	mux.HandleFunc("GET /api/v1/auth/refresh", h.Refresh)
	mux.Handle("GET /api/v1/auth/me", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.Me))))
	mux.HandleFunc("GET /api/v1/auth/logout", h.Logout)
}

func RegisterWebsiteRoutes(mux *http.ServeMux, h *handlers.WebSiteHandler) {
	mux.HandleFunc("POST /api/v1/websites", h.CreateWebSite)
}
