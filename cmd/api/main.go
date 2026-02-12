package main

import (
	"log"
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/config"
	"github.com/ViitoJooj/Jesterx/internal/database"
	httpRouter "github.com/ViitoJooj/Jesterx/internal/http"
	"github.com/ViitoJooj/Jesterx/internal/http/handlers"
	"github.com/ViitoJooj/Jesterx/internal/repository/postgres"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

func main() {
	config.LoadEnv()
	cfg := database.PostgresConfig(*config.PGCNN)

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	router := httpRouter.NewRouter(authHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
