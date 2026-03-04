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

	// Services
	authService := service.NewAuthService(authRepo, websiteRepo)
	websiteService := service.NewWebSiteService(websiteRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	websiteHandler := handlers.NewWebSiteHandler(websiteService)

	// Routers
	httpRouter.RegisterAuthRoutes(mux, authHandler, authService)
	httpRouter.RegisterWebsiteRoutes(mux, websiteHandler)

	// Middlewares
	handler := middleware.IdentityMiddleware(authService)(middleware.CORS(mux))

	go jobs.StartCleanupUserWorker(authService)

	http.ListenAndServe(":8080", handler)
}
