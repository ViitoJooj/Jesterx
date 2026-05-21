package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/config"
	apphttp "github.com/ViitoJooj/Jesterx/internal/http"
	"github.com/ViitoJooj/Jesterx/internal/http/handlers"
	mw "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/jobs"
	"github.com/ViitoJooj/Jesterx/internal/repository/postgres"
	"github.com/ViitoJooj/Jesterx/internal/service"
	"github.com/ViitoJooj/Jesterx/pkg/database"
	"github.com/ViitoJooj/Jesterx/pkg/logger"
	"github.com/ViitoJooj/Jesterx/pkg/migrate"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	if err := migrate.Run(db, "migrations"); err != nil {
		log.Fatalf("migrations: %v", err)
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
	storeSocialHandler := handlers.NewStoreSocialHandler(storeSocialService)

	router := apphttp.NewRouter(
		cfg,
		authService,
		authHandler,
		websiteHandler,
		paymentHandler,
		productHandler,
		orderHandler,
		storageHandler,
		themeHandler,
		adminHandler,
		reportHandler,
		storeSocialHandler,
	)

	handler := mw.CORS(
		logger.Middleware(func(ctx context.Context) string {
			id, ok := mw.UserID(ctx)
			if !ok {
				return ""
			}
			return id
		})(
			router,
		),
	)

	go jobs.StartCleanupUserWorker(authService)
	go jobs.StartSalesDigestWorker(orderService, authRepo, websiteRepo)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	log.Printf("server listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}
