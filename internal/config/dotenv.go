package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	IsDev       bool

	PGHost     string
	PGPort     string
	PGUser     string
	PGPassword string
	PGDBName   string
	PGSSLMode  string

	JWTAccessSecret  string
	JWTRefreshSecret string

	ResendKey           string
	StripePublic        string
	StripeSecret        string
	StripeWebhookSecret string

	FrontendURL string
	BackendURL  string
	StoragePath string

	PlatformCommissionPct string
}

var (
	PGCNN = &PostgresConnection{}

	Jwt_access_token  string
	Jwt_refresh_token string

	IsDev               bool
	ResendKey           string
	StripePublic        string
	StripeSecret        string
	StripeWebhookSecret string
	FrontendURL         string
	BackendURL          string
	StoragePath         string
)

type PostgresConnection struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Load reads environment variables and returns a Config struct.
// It also sets the package-level globals for backward compatibility.
func Load() *Config {
	_ = godotenv.Load(".env")

	cfg := &Config{
		Port:        getEnvOrDefault("PORT", "8080"),
		Environment: getEnvOrDefault("ENVIRONMENT", "dev"),

		PGHost:     mustGetenv("POSTGRES_HOST"),
		PGPort:     mustGetenv("POSTGRES_PORT"),
		PGUser:     mustGetenv("POSTGRES_USER"),
		PGPassword: mustGetenv("POSTGRES_PASSWORD"),
		PGDBName:   mustGetenv("POSTGRES_DB"),
		PGSSLMode:  getEnvOrDefault("POSTGRES_SSL", "disable"),

		JWTAccessSecret:  mustGetenv("JWT_ACCESS_TOKEN"),
		JWTRefreshSecret: mustGetenv("JWT_REFRESH_TOKEN"),

		ResendKey:           mustGetenv("RESEND_KEY"),
		StripePublic:        mustGetenv("STRIPE_PUBLIC_KEY"),
		StripeSecret:        mustGetenv("STRIPE_SECRET_KEY"),
		StripeWebhookSecret: getEnvOrDefault("STRIPE_WEBHOOK_SECRET", ""),

		FrontendURL: getEnvOrDefault("FRONTEND_URL", "http://localhost:5173"),
		BackendURL:  getEnvOrDefault("BACKEND_URL", "http://localhost:8080"),
		StoragePath: getEnvOrDefault("STORAGE_PATH", "./data"),

		PlatformCommissionPct: getEnvOrDefault("PLATFORM_COMMISSION_PCT", "5"),
	}
	cfg.IsDev = cfg.Environment == "dev"

	PGCNN.Host = cfg.PGHost
	PGCNN.Port = cfg.PGPort
	PGCNN.User = cfg.PGUser
	PGCNN.Password = cfg.PGPassword
	PGCNN.DBName = cfg.PGDBName
	PGCNN.SSLMode = cfg.PGSSLMode

	Jwt_access_token = cfg.JWTAccessSecret
	Jwt_refresh_token = cfg.JWTRefreshSecret

	IsDev = cfg.IsDev
	ResendKey = cfg.ResendKey
	StripePublic = cfg.StripePublic
	StripeSecret = cfg.StripeSecret
	StripeWebhookSecret = cfg.StripeWebhookSecret
	FrontendURL = cfg.FrontendURL
	BackendURL = cfg.BackendURL
	StoragePath = cfg.StoragePath

	return cfg
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var not set: %s", key)
	}
	return v
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
