package main

import (
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/middlewares"
	"gen-you-ecommerce/services"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	config.Load()
	if config.Gin_mode == "release" || config.Gin_mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	config.ConnectPostgres()
	config.ConnectMongo()
	router := gin.Default()

	// V1/AUTH
	router.POST("/v1/auth/login", middlewares.OptionalTenantMiddleware(), services.LoginService)
	router.POST("/v1/auth/register", middlewares.OptionalTenantMiddleware(), services.RegisterService)
	router.GET("/v1/auth/me", middlewares.OptionalTenantMiddleware(), middlewares.AuthMiddleware(), services.MeService)
	router.GET("/v1/auth/logout", middlewares.OptionalTenantMiddleware(), middlewares.AuthMiddleware(), services.LogoutService)

	// V1/PAGES
	router.POST("/v1/sites", middlewares.AuthMiddleware(), services.CreateSiteService)
	router.PUT("/v1/pages/:page_id", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.UpdatePageService)
	router.POST("/v1/pages", middlewares.TenantMiddleware(), middlewares.AuthMiddleware(), services.CreatePageService)
	router.GET("/v1/pages/:page_id", services.GetPageService)
	router.GET("/v1/pages/:page_id/raw", services.GetRawSveltePageService)

	http.ListenAndServe(":8080", router)
}
