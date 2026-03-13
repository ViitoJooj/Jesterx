package http

import (
	"net/http"
	"os"

	"github.com/ViitoJooj/Jesterx/internal/config"
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
	mux.Handle("DELETE /api/v1/auth/me", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.DeleteAccount))))
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

func RegisterStoreSocialRoutes(mux *http.ServeMux, h *handlers.StoreSocialHandler, authService *service.AuthService) {
	// Public
	mux.HandleFunc("GET /api/store/{siteID}/info", h.GetStoreFullInfo)
	mux.HandleFunc("GET /api/store/{siteID}/visits", h.GetVisitStats)
	mux.HandleFunc("GET /api/store/{siteID}/comments", h.ListComments)

	// Auth required
	mux.Handle("POST /api/store/{siteID}/comments",
		middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.PostComment))))
	mux.Handle("DELETE /api/store/{siteID}/comments/{commentID}",
		middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.DeleteComment))))
	mux.Handle("POST /api/store/{siteID}/ratings",
		middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.RateStore))))
	mux.Handle("GET /api/store/{siteID}/my-rating",
		middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.GetMyRating))))

	// Owner
	mux.Handle("PATCH /api/v1/sites/{siteID}/profile",
		middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.UpdateStoreProfile))))

	// Admin
	mux.Handle("PATCH /api/v1/admin/sites/{siteID}/mature",
		middleware.IdentityMiddleware(authService)(middleware.RequireRole(authService, "admin")(http.HandlerFunc(h.AdminSetMature))))

	// Team members (owner/manager/admin)
	mux.Handle("GET /api/v1/sites/{siteID}/members",
		middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.ListMembers))))
	mux.Handle("POST /api/v1/sites/{siteID}/members",
		middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.AddMember))))
	mux.Handle("PATCH /api/v1/sites/{siteID}/members/{memberUserID}",
		middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.UpdateMemberRole))))
	mux.Handle("DELETE /api/v1/sites/{siteID}/members/{memberUserID}",
		middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.RemoveMember))))

	// Comment replies (owner/manager/support/admin)
	mux.Handle("POST /api/store/{siteID}/comments/{commentID}/replies",
		middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.ReplyComment))))

	// My role in store (authenticated)
	mux.Handle("GET /api/store/{siteID}/my-role",
		middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.GetMyRole))))
}

func RegisterReportRoutes(mux *http.ServeMux, h *handlers.ReportHandler, authService *service.AuthService) {
	requireAdmin := middleware.RequireRole(authService, "admin")
	mux.HandleFunc("POST /api/v1/reports", h.PublicCreateReport)
	mux.Handle("GET /api/v1/admin/reports", middleware.IdentityMiddleware(authService)(requireAdmin(http.HandlerFunc(h.AdminListReports))))
	mux.Handle("GET /api/v1/admin/reports/{reportID}", middleware.IdentityMiddleware(authService)(requireAdmin(http.HandlerFunc(h.AdminGetReport))))
	mux.Handle("PATCH /api/v1/admin/reports/{reportID}", middleware.IdentityMiddleware(authService)(requireAdmin(http.HandlerFunc(h.AdminUpdateReport))))
}

func RegisterProductRoutes(mux *http.ServeMux, h *handlers.ProductHandler, authService *service.AuthService) {
	mux.Handle("POST /api/v1/sites/{siteID}/products", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.CreateProduct))))
	mux.Handle("GET /api/v1/sites/{siteID}/products", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.ListProducts))))
	mux.Handle("PATCH /api/v1/sites/{siteID}/products/{productID}", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.UpdateProduct))))
	mux.Handle("DELETE /api/v1/sites/{siteID}/products/{productID}", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.DeleteProduct))))
	mux.HandleFunc("GET /api/store/{siteID}/products", h.PublicListProducts)
	mux.HandleFunc("GET /api/store/{siteID}/products/{productID}", h.PublicGetProduct)
}

func RegisterOrderRoutes(mux *http.ServeMux, h *handlers.OrderHandler, auth *service.AuthService) {
	mux.Handle("POST /api/store/{siteID}/orders",
		middleware.IdentityMiddleware(auth)(middleware.RequireAuth(http.HandlerFunc(h.CreateOrder))))
	mux.Handle("GET /api/v1/sites/{siteID}/orders",
		middleware.IdentityMiddleware(auth)(middleware.RequireAuth(http.HandlerFunc(h.ListSiteOrders))))
}

func RegisterStorageRoutes(mux *http.ServeMux, h *handlers.StorageHandler, authService *service.AuthService) {
	mux.Handle("POST /api/v1/upload",
		middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.Upload))))

	// Serve uploaded files from the data directory.
	_ = os.MkdirAll(config.StoragePath, 0755)
	fileServer := http.FileServer(http.Dir(config.StoragePath))
	mux.Handle("GET /files/", http.StripPrefix("/files/", fileServer))
}

func RegisterThemeRoutes(mux *http.ServeMux, h *handlers.ThemeHandler) {
	mux.HandleFunc("GET /api/v1/themes", h.ListThemes)
}

func RegisterAdminRoutes(mux *http.ServeMux, h *handlers.AdminHandler, authService *service.AuthService) {
	requireAdmin := middleware.RequireRole(authService, "admin")
	mux.Handle("GET /api/v1/admin/stats", middleware.IdentityMiddleware(authService)(requireAdmin(http.HandlerFunc(h.Stats))))
	mux.Handle("GET /api/v1/admin/users", middleware.IdentityMiddleware(authService)(requireAdmin(http.HandlerFunc(h.ListUsers))))
	mux.Handle("GET /api/v1/admin/sites", middleware.IdentityMiddleware(authService)(requireAdmin(http.HandlerFunc(h.ListSites))))
	mux.Handle("GET /api/v1/admin/orders", middleware.IdentityMiddleware(authService)(requireAdmin(http.HandlerFunc(h.ListOrders))))
	mux.Handle("GET /api/v1/admin/revenue", middleware.IdentityMiddleware(authService)(requireAdmin(http.HandlerFunc(h.Revenue))))
}

func RegisterPaymentRoutes(mux *http.ServeMux, h *handlers.PaymentHandler, authService *service.AuthService) {
	mux.HandleFunc("GET /api/v1/plans", h.ListPlans)
	mux.Handle("POST /api/v1/payments/checkout", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.CreateCheckout))))
	mux.Handle("GET /api/v1/payments/confirm", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.ConfirmCheckout))))
	mux.Handle("POST /api/v1/payments/cancel", middleware.IdentityMiddleware(authService)(middleware.RequireAuth(http.HandlerFunc(h.CancelSubscription))))
	mux.HandleFunc("POST /api/v1/payments/webhook", h.StripeWebhook)
}
