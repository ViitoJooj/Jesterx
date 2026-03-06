package main

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/config"
	httpRouter "github.com/ViitoJooj/Jesterx/internal/http"
	"github.com/ViitoJooj/Jesterx/internal/http/handlers"
	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/jobs"
	"github.com/ViitoJooj/Jesterx/internal/repository/postgres"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

func main() {
	config.LoadEnv()
	mux := httpRouter.NewRouter()
	db := postgres.NewPostgres(postgres.PostgresConfig(*config.PGCNN))

	// Repositorys
	authRepo := postgres.NewAuthRepository(db)
	websiteRepo := postgres.NewWebSiteRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	productRepo := postgres.NewProductRepository(db)

	// Services
	authService := service.NewAuthService(authRepo, websiteRepo, paymentRepo)
	websiteService := service.NewWebSiteService(websiteRepo, authRepo, paymentRepo)
	paymentService := service.NewPaymentService(paymentRepo, authRepo)
	productService := service.NewProductService(productRepo, websiteRepo, authRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	websiteHandler := handlers.NewWebSiteHandler(websiteService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	productHandler := handlers.NewProductHandler(productService)

	// Routers
	httpRouter.RegisterAuthRoutes(mux, authHandler, authService)
	httpRouter.RegisterWebsiteRoutes(mux, websiteHandler, authService)
	httpRouter.RegisterPaymentRoutes(mux, paymentHandler, authService)
	httpRouter.RegisterProductRoutes(mux, productHandler, authService)

	// Middlewares
	handler := middleware.IdentityMiddleware(authService)(middleware.CORS(mux))

	go jobs.StartCleanupUserWorker(authService)

	http.ListenAndServe(":8080", handler)
}
