package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	JwtSecret           string
	PostgresUri         string
	MongoUri            string
	GinMode             string
	StripeSecretKey     string
	StripeWebhookSecret string
)

func Load() {
	_ = godotenv.Load("../.env")

	postgresUser := mustGetenv("POSTGRES_USER")
	postgresPassword := mustGetenv("POSTGRES_PASSWORD")
	postgresDB := mustGetenv("POSTGRES_DB")
	postgresPort := mustGetenv("POSTGRES_PORT")
	postgresHost := mustGetenv("POSTGRES_HOST")
	postgresSSL := mustGetenv("POSTGRES_SSL")

	PostgresUri = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB, postgresSSL,
	)

	mongoUser := mustGetenv("MONGO_USER")
	mongoPassword := mustGetenv("MONGO_PASSWORD")
	mongoPort := mustGetenv("MONGO_PORT")
	mongoHost := mustGetenv("MONGO_HOST")

	MongoUri = fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUser, mongoPassword, mongoHost, mongoPort)

	JwtSecret = mustGetenv("JWT_SECRET")
	GinMode = mustGetenv("GIN_MODE")

	StripeSecretKey = mustGetenv("STRIPE_API_SECRET")
	StripeWebhookSecret = mustGetenv("STRIPE_WEBHOOK_SECRET")
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatal("Error on get " + key)
	}
	return v
}
