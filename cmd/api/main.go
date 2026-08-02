package main

import (
	"net/http"
	"os"

	"github.com/ViitoJooj/Jesterx/internal/controllers"
	"github.com/ViitoJooj/Jesterx/internal/repositories"
	"github.com/ViitoJooj/Jesterx/internal/usecases"
	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	"github.com/ViitoJooj/Jesterx/pkg/logger"
	"github.com/ViitoJooj/Jesterx/pkg/postgresql"
)

func main() {
	config, err := dotenv.Conn()
	if err != nil {
		logger.Warn(err).Print()
		os.Exit(0)
	}

	db, err := postgresql.Conn(config.PostgreSQL)
	if err != nil {
		logger.Warn(err).Print()
		os.Exit(0)
	}

	err = postgresql.Migrations(db)
	if err != nil {
		logger.Warn(err).Print()
		os.Exit(0)
	}

	userRepo := repositories.NewUserRepository(db)

	registerUseCase := usecases.NewRegisterUserUseCase(db, userRepo)

	authController := controllers.NewAuthController(registerUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/register", authController.Register)

	http.ListenAndServe(":8080", mux)
}
