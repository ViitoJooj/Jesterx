package main

import (
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

	// V1/BILLING
	router.POST("/v1/billing/checkout", middlewares.AuthMiddleware(), services.CreateCheckoutService)
	router.POST("/v1/billing/webhook", services.PaymentWebhookService)
	router.POST("/v1/billing/confirm", middlewares.AuthMiddleware(), services.ConfirmCheckoutService)

	fmt.Print("\nYour api is running:\n\n")
	fmt.Println("Bind address:", "0.0.0.0:"+config.ApplicationPort)
	fmt.Println("LAN (same network): " + ip + ":" + config.ApplicationPort)
	fmt.Println("Local (this machine): " + "http://localhost:" + config.ApplicationPort)
	fmt.Println("Website: " + config.HostProd)

	if err := http.ListenAndServe(":"+config.ApplicationPort, router); err != nil {
		panic(err)
	}
}
