package main

import (
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/middlewares"
	"gen-you-ecommerce/services"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	config.Load()
	if config.GinMode == "release" || config.GinMode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	config.ConnectPostgres()
	config.ConnectMongo()
	config.InitStripe()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Tenant-Page-Id"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// V1/AUTH
	router.POST("/v1/auth/login", middlewares.OptionalTenantMiddleware(), services.LoginService)
	router.POST("/v1/auth/register", middlewares.OptionalTenantMiddleware(), services.RegisterService)
	router.GET("/v1/auth/me", middlewares.OptionalTenantMiddleware(), middlewares.AuthMiddleware(), services.MeService)
	router.GET("/v1/auth/logout", middlewares.OptionalTenantMiddleware(), middlewares.AuthMiddleware(), services.LogoutService)
	router.GET("/v1/auth/refresh", middlewares.OptionalTenantMiddleware(), middlewares.AuthMiddleware(), services.RefreshUserService)

	// V1/SITES
	router.POST("/v1/sites", middlewares.AuthMiddleware(), middlewares.PlanMiddleware(), services.CreateSiteService)

	// V1/PAGES
	router.PUT("/v1/pages/:page_id", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.UpdatePageService)
	router.GET("/v1/pages", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.ListPagesService)
	router.POST("/v1/pages", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.CreatePageService)
	router.GET("/v1/pages/:page_id", services.GetPageService)
	router.GET("/v1/pages/:page_id/raw", services.GetRawSveltePageService)

	// V1/BILLING
	router.POST("/v1/billing/checkout", middlewares.AuthMiddleware(), services.CreateCheckoutService)
	router.POST("/v1/billing/webhook", services.PaymentWebhookService)
	router.POST("/v1/billing/confirm", middlewares.AuthMiddleware(), services.ConfirmCheckoutService)

	http.ListenAndServe(":8080", router)
}
