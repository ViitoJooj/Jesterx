package main

import (
	"log"

	httpx "github.com/ViitoJooj/Jesterx/internal/http"
	"github.com/ViitoJooj/Jesterx/internal/http/handlers"
	"github.com/ViitoJooj/Jesterx/internal/repository"
	"github.com/ViitoJooj/Jesterx/internal/service"
	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	"github.com/ViitoJooj/Jesterx/pkg/redis"
	"github.com/ViitoJooj/Jesterx/pkg/supabase"
	"github.com/ViitoJooj/Jesterx/pkg/validators"
)

func main() {
	dotenv.Conn()
	supabase.Conn()
	redis.Conn()
	validators.LoadEmbedded()

	userRepo := repository.NewUserRepository(supabase.DB)
	authService := service.NewAuthService(userRepo, redis.Client)
	authHandler := handlers.NewAuthHandler(authService)

	router := httpx.NewRouter(authHandler)

	log.Println("Running on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Panicf("Error running server: %s", err)
	}
}
