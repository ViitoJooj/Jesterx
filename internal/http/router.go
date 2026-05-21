package http

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/ViitoJooj/Jesterx/internal/config"
	"github.com/ViitoJooj/Jesterx/internal/http/handlers"
	mw "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

// NewRouter builds the chi mux with all route registrations.
// Identity extraction is applied globally so any route can optionally read
// the authenticated user from context, whether or not it requires auth.
func NewRouter(
	cfg *config.Config,
	authService *service.AuthService,
	authHandler *handlers.AuthHandler,
	websiteHandler *handlers.WebSiteHandler,
	paymentHandler *handlers.PaymentHandler,
	productHandler *handlers.ProductHandler,
	orderHandler *handlers.OrderHandler,
	storageHandler *handlers.StorageHandler,
	themeHandler *handlers.ThemeHandler,
	adminHandler *handlers.AdminHandler,
	reportHandler *handlers.ReportHandler,
	storeSocialHandler *handlers.StoreSocialHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(mw.IdentityMiddleware(authService))

	_ = os.MkdirAll(cfg.StoragePath, 0755)
	fileServer := http.FileServer(http.Dir(cfg.StoragePath))
	r.Handle("/files/*", http.StripPrefix("/files/", fileServer))

	r.Get("/p/{siteID}", websiteHandler.PublicRender)
	r.Get("/p/{siteID}/*", websiteHandler.PublicRender)

	r.Route("/api/store/{siteID}", func(r chi.Router) {
		r.Get("/info", storeSocialHandler.GetStoreFullInfo)
		r.Get("/visits", storeSocialHandler.GetVisitStats)
		r.Get("/comments", storeSocialHandler.ListComments)
		r.Get("/products", productHandler.PublicListProducts)
		r.Get("/products/{productID}", productHandler.PublicGetProduct)

		r.Group(func(r chi.Router) {
			r.Use(mw.RequireAuth)
			r.Post("/comments", storeSocialHandler.PostComment)
			r.Delete("/comments/{commentID}", storeSocialHandler.DeleteComment)
			r.Post("/comments/{commentID}/replies", storeSocialHandler.ReplyComment)
			r.Post("/ratings", storeSocialHandler.RateStore)
			r.Get("/my-rating", storeSocialHandler.GetMyRating)
			r.Get("/my-role", storeSocialHandler.GetMyRole)
			r.Post("/orders", orderHandler.CreateOrder)
		})
	})

	r.Route("/api/v1", func(r chi.Router) {

		r.Get("/plans", paymentHandler.ListPlans)
		r.Get("/themes", themeHandler.ListThemes)
		r.Post("/reports", reportHandler.PublicCreateReport)
		r.Post("/payments/webhook", paymentHandler.StripeWebhook)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Get("/verify/{id}", authHandler.VerifyEmail)
			r.Get("/refresh", authHandler.Refresh)
			r.Get("/logout", authHandler.Logout)

			r.Group(func(r chi.Router) {
				r.Use(mw.RequireAuth)
				r.Get("/me", authHandler.Me)
				r.Patch("/me", authHandler.UpdateProfile)
				r.Delete("/me", authHandler.DeleteAccount)

				r.Get("/addresses", authHandler.ListAddresses)
				r.Post("/addresses", authHandler.CreateAddress)
				r.Patch("/addresses/{id}", authHandler.UpdateAddress)
				r.Delete("/addresses/{id}", authHandler.DeleteAddress)
				r.Post("/addresses/{id}/default", authHandler.SetDefaultAddress)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(mw.RequireAuth)

			r.Get("/websites", websiteHandler.ListWebSites)
			r.Get("/site-apis", websiteHandler.ListSiteAPIs)
			r.Post("/websites", websiteHandler.CreateWebSite)

			r.Route("/sites/{siteID}", func(r chi.Router) {
				r.Delete("/", websiteHandler.DeleteWebSite)
				r.Get("/routes", websiteHandler.ListRoutes)
				r.Post("/routes", websiteHandler.ReplaceRoutes)
				r.Get("/versions", websiteHandler.ListVersions)
				r.Post("/versions", websiteHandler.CreateVersion)
				r.Post("/publish/{version}", websiteHandler.PublishVersion)
				r.Get("/scan-reports/{version}", websiteHandler.GetScanReport)
				r.Patch("/profile", storeSocialHandler.UpdateStoreProfile)

				r.Get("/products", productHandler.ListProducts)
				r.Post("/products", productHandler.CreateProduct)
				r.Patch("/products/{productID}", productHandler.UpdateProduct)
				r.Delete("/products/{productID}", productHandler.DeleteProduct)

				r.Get("/orders", orderHandler.ListSiteOrders)

				r.Get("/members", storeSocialHandler.ListMembers)
				r.Post("/members", storeSocialHandler.AddMember)
				r.Patch("/members/{memberUserID}", storeSocialHandler.UpdateMemberRole)
				r.Delete("/members/{memberUserID}", storeSocialHandler.RemoveMember)
			})

			r.Post("/payments/checkout", paymentHandler.CreateCheckout)
			r.Get("/payments/confirm", paymentHandler.ConfirmCheckout)
			r.Post("/payments/cancel", paymentHandler.CancelSubscription)

			r.Post("/upload", storageHandler.Upload)
		})

		r.Group(func(r chi.Router) {
			r.Use(mw.RequireRole(authService, "admin"))
			r.Get("/admin/stats", adminHandler.Stats)
			r.Get("/admin/users", adminHandler.ListUsers)
			r.Get("/admin/sites", adminHandler.ListSites)
			r.Get("/admin/orders", adminHandler.ListOrders)
			r.Get("/admin/revenue", adminHandler.Revenue)
			r.Get("/admin/reports", reportHandler.AdminListReports)
			r.Get("/admin/reports/{reportID}", reportHandler.AdminGetReport)
			r.Patch("/admin/reports/{reportID}", reportHandler.AdminUpdateReport)
			r.Patch("/admin/sites/{siteID}/mature", storeSocialHandler.AdminSetMature)
		})
	})

	return r
}
