package main

import (
	"context"
	"fmt"
	"jesterx-core/config"
	"jesterx-core/middlewares"
	"jesterx-core/services"
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

	ip, err := config.HostIP()
	if err != nil {
		panic(err)
	}

	config.ConnectPostgres()
	config.ConnectMongo()
	config.InitStripe()
	config.InitOAuth()
	if err := services.SetupPlatformData(context.Background()); err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
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

	// V1/AUTH/OAUTH2
	router.GET("/v1/auth/google", services.GoogleLoginService)
	router.GET("/v1/auth/google/callback", services.GoogleCallbackService)
	router.GET("/v1/auth/github", services.GithubLoginService)
	router.GET("/v1/auth/github/callback", services.GithubCallbackService)
	router.GET("/v1/auth/twitter", services.TwitterLoginService)
	router.GET("/v1/auth/twitter/callback", services.TwitterCallbackService)

	// V1/SITES
	router.POST("/v1/sites", middlewares.AuthMiddleware(), middlewares.PlanMiddleware(), services.CreateSiteService)

	// V1/PAGES
	router.PUT("/v1/pages/:page_id", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.UpdatePageService)
	router.GET("/v1/pages", middlewares.AuthMiddleware(), middlewares.OptionalTenantMiddleware(), services.ListPagesService)
	router.POST("/v1/pages", middlewares.AuthMiddleware(), middlewares.OptionalTenantMiddleware(), services.CreatePageService)
	router.GET("/v1/pages/:page_id", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.GetPageService)
	router.GET("/v1/pages/:page_id/raw", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.GetRawSveltePageService)
	router.DELETE("/v1/pages/:page_id", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.DeletePageService)
	router.GET("/v1/pages/:page_id/products", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.ListProductsService)
	router.POST("/v1/pages/:page_id/products", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.CreateProductService)
	router.PUT("/v1/pages/:page_id/products/:product_id", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.UpdateProductService)

	// Public Plans
	router.GET("/v1/plans", services.ListPlansService)

	// V1/THEMES
	router.POST("/v1/themes/apply", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.ApplyThemeService)
	router.GET("/v1/themes/store", middlewares.OptionalTenantMiddleware(), services.ListThemeStoreService)
	router.GET("/v1/themes/store/:slug", services.GetThemeStoreBySlugService)
	router.PUT("/v1/themes/store/:page_id", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.UpdateThemeStoreEntryService)

	// V1/BILLING
	router.POST("/v1/billing/checkout", middlewares.AuthMiddleware(), services.CreateCheckoutService)
	router.POST("/v1/billing/webhook", services.PaymentWebhookService)
	router.POST("/v1/billing/confirm", middlewares.AuthMiddleware(), services.ConfirmCheckoutService)

	// V1/ADMIN
	admin := router.Group("/v1/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		admin.GET("/plans", services.AdminListPlansService)
		admin.PUT("/plans/:plan_id", services.AdminUpdatePlanService)
		admin.GET("/users", services.AdminListUsersService)
		admin.PUT("/users/:user_id", services.AdminUpdateUserService)
		admin.PUT("/users/:user_id/ban", services.AdminBanUserService)
		admin.DELETE("/users/:user_id", services.AdminDeleteUserService)
		admin.GET("/users/export", services.AdminExportUsersService)
		admin.GET("/stats/overview", services.AdminOverviewService)
	}

	// Public access to page content and products
	router.GET("/v1/public/pages/:page_id", middlewares.TenantMiddleware(), services.PublicPageService)
	router.GET("/v1/public/pages/:page_id/products", middlewares.TenantMiddleware(), services.PublicListProductsService)

	fmt.Print("\nYour api is running:\n\n")
	fmt.Println("Bind address:", "0.0.0.0:"+config.ApplicationPort)
	fmt.Println("LAN (same network): " + ip + ":" + config.ApplicationPort)
	fmt.Println("Local (this machine): " + "http://localhost:" + config.ApplicationPort)
	fmt.Println("Website: " + config.HostProd)

	if err := http.ListenAndServe(":"+config.ApplicationPort, router); err != nil {
		panic(err)
	}
}
