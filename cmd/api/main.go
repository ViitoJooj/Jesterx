package main

import (
	"log"

	httpx "github.com/ViitoJooj/Jesterx/internal/http"
	"github.com/ViitoJooj/Jesterx/internal/repository"
	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	"github.com/ViitoJooj/Jesterx/pkg/supabase"
	"github.com/ViitoJooj/Jesterx/pkg/validators"
)

func main() {
	dotenv.Set()
	supabase.Conn()
	validators.LoadEmbedded()

	router := httpx.RegisterRouters()

	userRepository := repository.NewUserRepository(supabase.DB)
	_ = userRepository

	log.Println("Running...")
	if err := router.Run(); err != nil {
		log.Panicf("Error running server: %s", err)
	}
}
