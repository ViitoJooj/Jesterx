package main

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/config"
	httpRouter "github.com/ViitoJooj/Jesterx/internal/http"
	"github.com/ViitoJooj/Jesterx/internal/http/handlers"
	middleware "github.com/ViitoJooj/Jesterx/internal/http/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/repository/postgres"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

func main() {
	config.LoadEnv()
	mux := httpRouter.NewRouter()
	db := postgres.NewPostgres(postgres.PostgresConfig(*config.PGCNN))
	repo := postgres.NewRepository(db)

	authService := service.NewAuthService(repo)
	websiteService := service.NewWebSiteService(repo)

	authHandler := handlers.NewAuthHandler(authService)
	websiteHandler := handlers.NewWebSiteHandler(websiteService)

	httpRouter.RegisterAuthRoutes(mux, authHandler)
	httpRouter.RegisterWebsiteRoutes(mux, websiteHandler)

	handler := middleware.CORS(mux)
	http.ListenAndServe(":8080", handler)
}
