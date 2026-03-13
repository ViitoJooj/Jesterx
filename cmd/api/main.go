package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	"github.com/ViitoJooj/Jesterx/pkg/safeguard"
)

const (
	maxBodyBytes   = 10 * 1024 * 1024
	maxUploadBytes = 50 * 1024 * 1024
	maxPaginationN = 100
	banStrikes     = 20
	banDuration    = 30 * time.Minute
)

func main() {
	config.LoadEnv()
	mux := httpRouter.NewRouter()
	db := postgres.NewPostgres(postgres.PostgresConfig(*config.PGCNN))

	if err := migrate.Run(db, "migrations"); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	authRepo := postgres.NewAuthRepository(db)
	websiteRepo := postgres.NewWebSiteRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	productRepo := postgres.NewProductRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	reportRepo := postgres.NewReportRepository(db)
	storeSocialRepo := postgres.NewStoreSocialRepository(db)

	authService := service.NewAuthService(authRepo, websiteRepo, paymentRepo)
	websiteService := service.NewWebSiteService(websiteRepo, authRepo, paymentRepo)
	paymentService := service.NewPaymentService(paymentRepo, authRepo)
	productService := service.NewProductService(productRepo, websiteRepo, authRepo)
	orderService := service.NewOrderService(orderRepo, websiteRepo, productRepo, authRepo)
	reportService := service.NewReportService(reportRepo, websiteRepo)
	storeSocialService := service.NewStoreSocialService(storeSocialRepo, websiteRepo)

	storageService := service.NewStorageService()

	authHandler := handlers.NewAuthHandler(authService)
	websiteHandler := handlers.NewWebSiteHandler(websiteService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)
	storageHandler := handlers.NewStorageHandler(storageService)
	themeHandler := handlers.NewThemeHandler(db)
	adminHandler := handlers.NewAdminHandler(db)
	reportHandler := handlers.NewReportHandler(reportService, authService)
	storeSocialHandler := handlers.NewStoreSocialHandler(storeSocialService, db)

	httpRouter.RegisterAuthRoutes(mux, authHandler, authService)
	httpRouter.RegisterWebsiteRoutes(mux, websiteHandler, authService)
	httpRouter.RegisterPaymentRoutes(mux, paymentHandler, authService)
	httpRouter.RegisterProductRoutes(mux, productHandler, authService)
	httpRouter.RegisterOrderRoutes(mux, orderHandler, authService)
	httpRouter.RegisterStorageRoutes(mux, storageHandler, authService)
	httpRouter.RegisterThemeRoutes(mux, themeHandler)
	httpRouter.RegisterAdminRoutes(mux, adminHandler, authService)
	httpRouter.RegisterReportRoutes(mux, reportHandler, authService)
	httpRouter.RegisterStoreSocialRoutes(mux, storeSocialHandler, authService)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
			log.Printf("health encode error: %v", err)
		}
	})

	globalLimiter := ratelimit.NewLimiter(200)
	authLimiter := ratelimit.NewLimiter(15)

	routeLimiter := ratelimit.NewRouteRateLimiter().
		Add("/api/v1/payments/", 10).
		Add("/api/v1/upload", 20).
		Add("/api/store/", 120)

	ipBanner := safeguard.NewIPBanner(banStrikes, banDuration)

	handler := safeguard.Recovery(
		logger.Middleware(func(ctx context.Context) string {
			id, ok := middleware.UserID(ctx)
			if !ok {
				return ""
			}
			return id
		})(
			middleware.CORS(
				safeguard.PathTraversalGuard(
					safeguard.PaginationGuard(maxPaginationN)(
						safeguard.BodyLimit(maxBodyBytes)(
							ipBanner.Middleware(
								globalLimiter.Middleware(
									routeLimiter.Middleware(
										ratelimit.AuthRateLimit(authLimiter,
											middleware.IdentityMiddleware(authService)(mux),
										),
									),
								),
							),
						),
					),
				),
			),
		),
	)

	go jobs.StartCleanupUserWorker(authService)
	go jobs.StartSalesDigestWorker(orderService, authRepo, websiteRepo)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
