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
	db := postgres.NewPostgres(postgres.PostgresConfig(*config.PGCNN))
	userRepo := postgres.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)
	router := httpRouter.Routers(authHandler)
	handler := middleware.CORS(router)

	http.ListenAndServe(":8080", handler)
}
