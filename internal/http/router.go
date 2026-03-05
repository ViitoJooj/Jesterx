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
	mux.Handle("PATCH /api/v1/auth/me", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.UpdateProfile))))
	mux.HandleFunc("GET /api/v1/auth/logout", h.Logout)
}

func RegisterWebsiteRoutes(mux *http.ServeMux, h *handlers.WebSiteHandler, authService *service.AuthService) {
	mux.Handle("GET /p/{siteID}/{path...}", http.HandlerFunc(h.PublicRender))
	mux.Handle("GET /p/{siteID}", http.HandlerFunc(h.PublicRender))

	mux.Handle("GET /api/v1/websites", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.ListWebSites))))
	mux.Handle("GET /api/v1/site-apis", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.ListSiteAPIs))))
	mux.Handle("POST /api/v1/websites", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.CreateWebSite))))
	mux.Handle("DELETE /api/v1/sites/{siteID}", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.DeleteWebSite))))
	mux.Handle("POST /api/v1/sites/{siteID}/routes", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.ReplaceRoutes))))
	mux.Handle("GET /api/v1/sites/{siteID}/routes", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.ListRoutes))))
	mux.Handle("GET /api/v1/sites/{siteID}/versions", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.ListVersions))))
	mux.Handle("POST /api/v1/sites/{siteID}/versions", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.CreateVersion))))
	mux.Handle("POST /api/v1/sites/{siteID}/publish/{version}", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.PublishVersion))))
	mux.Handle("GET /api/v1/sites/{siteID}/scan-reports/{version}", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.GetScanReport))))
}

func RegisterPaymentRoutes(mux *http.ServeMux, h *handlers.PaymentHandler, authService *service.AuthService) {
	mux.HandleFunc("GET /api/v1/plans", h.ListPlans)
	mux.Handle("POST /api/v1/payments/checkout", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.CreateCheckout))))
	mux.Handle("GET /api/v1/payments/confirm", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.ConfirmCheckout))))
	mux.Handle("POST /api/v1/payments/cancel", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.CancelSubscription))))
	mux.HandleFunc("POST /api/v1/payments/webhook", h.StripeWebhook)
}
