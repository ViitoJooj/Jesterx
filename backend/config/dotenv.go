package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	JwtSecret           string
	PostgresUri         string
	MongoUri            string
	GinMode             string
	StripeSecretKey     string
	StripeWebhookSecret string
	HostProd            string
	ApplicationPort     string
	AdminEmails         []string
)

func Load() {
	_ = godotenv.Load("../.env")

	//App
	HostProd = mustGetenv("HOST_PROD")
	JwtSecret = mustGetenv("JWT_SECRET")
	GinMode = mustGetenv("GIN_MODE")
	ApplicationPort = mustGetenv("APPLICATION_PORT")

	//Postgres
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

	adminList := os.Getenv("ADMIN_EMAILS")
	if adminList != "" {
		for _, email := range strings.Split(adminList, ",") {
			email = strings.TrimSpace(email)
			if email != "" {
				AdminEmails = append(AdminEmails, email)
			}
		}
	}

	//MongoDB
	mongoUser := mustGetenv("MONGO_USER")
	mongoPassword := mustGetenv("MONGO_PASSWORD")
	mongoPort := mustGetenv("MONGO_PORT")
	mongoHost := mustGetenv("MONGO_HOST")

	MongoUri = fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUser, mongoPassword, mongoHost, mongoPort)

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
