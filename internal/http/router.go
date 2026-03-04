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
	mux.HandleFunc("GET /api/v1/auth/verify/", h.VerifyEmail)
	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
	mux.HandleFunc("GET /api/v1/auth/refresh", h.Refresh)
	mux.Handle("GET /api/v1/auth/me", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.Me))))
	mux.HandleFunc("GET /api/v1/auth/logout", h.Logout)
}

func RegisterWebsiteRoutes(mux *http.ServeMux, h *handlers.WebSiteHandler) {
	mux.HandleFunc("POST /api/v1/websites", h.CreateWebSite)
}

func RegisterPaymentRoutes(mux *http.ServeMux, h *handlers.PaymentHandler, authService *service.AuthService) {
	mux.HandleFunc("GET /api/v1/plans", h.ListPlans)
	mux.Handle("POST /api/v1/payments/checkout", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.CreateCheckout))))
	mux.Handle("GET /api/v1/payments/confirm", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.ConfirmCheckout))))
	mux.HandleFunc("POST /api/v1/payments/webhook", h.StripeWebhook)
}
