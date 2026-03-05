package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type PostgresConnection struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var PGCNN = &PostgresConnection{}
var Jwt_access_token string
var Jwt_refresh_token string
var IsDev bool
var ResendKey string
var StripePublic string
var StripeSecret string
var StripeWebhookSecret string
var FrontendURL string
var SupabaseURL string
var SupabaseAnonKey string

func LoadEnv() {
	_ = godotenv.Load(".env")

	PGCNN.User = mustGetenv("POSTGRES_USER")
	PGCNN.Password = mustGetenv("POSTGRES_PASSWORD")
	PGCNN.DBName = mustGetenv("POSTGRES_DB")
	PGCNN.Port = mustGetenv("POSTGRES_PORT")
	PGCNN.Host = mustGetenv("POSTGRES_HOST")
	PGCNN.SSLMode = mustGetenv("POSTGRES_SSL")

	Jwt_access_token = mustGetenv("JWT_ACCESS_TOKEN")
	Jwt_refresh_token = mustGetenv("JWT_REFRESH_TOKEN")

	environment := mustGetenv("ENVIRONMENT")
	if environment == "dev" {
		IsDev = true
	}

	ResendKey = mustGetenv("RESEND_KEY")

	StripePublic = mustGetenv("STRIPE_PUBLIC_KEY")
	StripeSecret = mustGetenv("STRIPE_SECRET_KEY")
	StripeWebhookSecret = getEnvOrDefault("STRIPE_WEBHOOK_SECRET", "")
	FrontendURL = getEnvOrDefault("FRONTEND_URL", "http://localhost:5173")
	SupabaseURL = getEnvOrDefault("SUPABASE_PROJECT_URL", "")
	SupabaseAnonKey = getEnvOrDefault("SUPABASE_ANON_KEY", "")
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatal("Error on get " + key)
	}
	return v
}

func getEnvOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
