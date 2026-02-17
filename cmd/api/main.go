package main

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/config"
	"github.com/ViitoJooj/Jesterx/internal/database"
	httpRouter "github.com/ViitoJooj/Jesterx/internal/http"
	"github.com/ViitoJooj/Jesterx/internal/http/handlers"
	middleware "github.com/ViitoJooj/Jesterx/internal/middlewares"
	"github.com/ViitoJooj/Jesterx/internal/repository/postgres"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

func main() {
	config.LoadEnv()
	db := database.NewPostgres(database.PostgresConfig(*config.PGCNN))
	userRepo := postgres.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)
	router := httpRouter.AuthRouters(authHandler)
	handler := middleware.CORS(router)

	http.ListenAndServe(":8080", handler)
}
