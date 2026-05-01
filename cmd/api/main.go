package main

import (
	"log"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	postgres "github.com/ViitoJooj/Jesterx/pkg/postgres"
	"github.com/ViitoJooj/Jesterx/pkg/redis"
	validators "github.com/ViitoJooj/Jesterx/pkg/validators/users_validations"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := dotenv.Load()

	db, err := postgres.Conn(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redis.Conn(cfg.Redis)
	validators.LoadEmbedded()

	router := gin.Default()

	log.Println("Running on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Panicf("Error running server: %s", err)
	}
}
