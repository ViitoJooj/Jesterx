package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/config"
	httpRouter "github.com/ViitoJooj/Jesterx/internal/http"
	"github.com/ViitoJooj/Jesterx/internal/http/handlers"
	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/jobs"
	"github.com/ViitoJooj/Jesterx/internal/repository/postgres"
	"github.com/ViitoJooj/Jesterx/internal/service"
	"github.com/ViitoJooj/Jesterx/pkg/logger"
	"github.com/ViitoJooj/Jesterx/pkg/migrate"
	"github.com/ViitoJooj/Jesterx/pkg/ratelimit"
)

func main() {
	config.LoadEnv()
	mux := httpRouter.NewRouter()
	db := postgres.NewPostgres(postgres.PostgresConfig(*config.PGCNN))

	// Auto-run migrations
	if err := migrate.Run(db, "migrations"); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// Repositorys
	authRepo := postgres.NewAuthRepository(db)
	websiteRepo := postgres.NewWebSiteRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	productRepo := postgres.NewProductRepository(db)
	orderRepo := postgres.NewOrderRepository(db)

	// Services
	authService := service.NewAuthService(authRepo, websiteRepo, paymentRepo)
	websiteService := service.NewWebSiteService(websiteRepo, authRepo, paymentRepo)
	paymentService := service.NewPaymentService(paymentRepo, authRepo)
	productService := service.NewProductService(productRepo, websiteRepo, authRepo)
	orderService := service.NewOrderService(orderRepo, websiteRepo, productRepo)

	storageService := service.NewStorageService()

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	websiteHandler := handlers.NewWebSiteHandler(websiteService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)
	storageHandler := handlers.NewStorageHandler(storageService)
	themeHandler := handlers.NewThemeHandler(db)
	adminHandler := handlers.NewAdminHandler(db)

	// Routers
	httpRouter.RegisterAuthRoutes(mux, authHandler, authService)
	httpRouter.RegisterWebsiteRoutes(mux, websiteHandler, authService)
	httpRouter.RegisterPaymentRoutes(mux, paymentHandler, authService)
	httpRouter.RegisterProductRoutes(mux, productHandler, authService)
	httpRouter.RegisterOrderRoutes(mux, orderHandler, authService)
	httpRouter.RegisterStorageRoutes(mux, storageHandler, authService)
	httpRouter.RegisterThemeRoutes(mux, themeHandler)
	httpRouter.RegisterAdminRoutes(mux, adminHandler, authService)

	// Middlewares
	globalLimiter := ratelimit.NewLimiter(200)
	authLimiter := ratelimit.NewLimiter(15)

	handler := logger.Middleware(func(ctx context.Context) string {
		id, ok := middleware.UserID(ctx)
		if !ok {
			return ""
		}
		return id
	})(middleware.CORS(
		globalLimiter.Middleware(
			ratelimit.AuthRateLimit(authLimiter,
				middleware.IdentityMiddleware(authService)(mux),
			),
		),
	))

	go jobs.StartCleanupUserWorker(authService)
	go jobs.StartSalesDigestWorker(orderService, authRepo, websiteRepo)

	http.ListenAndServe(":8080", handler)
}
