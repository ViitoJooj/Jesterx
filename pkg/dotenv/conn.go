package dotenv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Postgres
type PostgresConfig struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName   string
	PostgresSSLMode  string
}

// Resend
type ResendConfig struct {
	ResendKey string
}

// JWT
type JWTConfig struct {
	SecretKey             string
	AccessTokenExpMinutes string
	RefreshTokenExpDays   string
}

// Redis
type RedisConfig struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
}

// Application
type ApplicationConfig struct {
	Environment string
}

type Config struct {
	Postgres    PostgresConfig
	Resend      ResendConfig
	JWT         JWTConfig
	Redis       RedisConfig
	Application ApplicationConfig
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Error loading .env: %v", err)
	}

	return &Config{
		Postgres: PostgresConfig{
			PostgresHost:     get("POSTGRES_HOST"),
			PostgresPort:     get("POSTGRES_PORT"),
			PostgresUser:     get("POSTGRES_USER"),
			PostgresPassword: get("POSTGRES_PASSWORD"),
			PostgresDBName:   get("POSTGRES_DBNAME"),
			PostgresSSLMode:  get("POSTGRES_SSLMODE"),
		},

		Resend: ResendConfig{
			ResendKey: get("RESEND_API_KEY"),
		},

		JWT: JWTConfig{
			SecretKey:             get("SECRET_KEY"),
			AccessTokenExpMinutes: get("ACCESS_TOKEN_EXP_MINUTES"),
			RefreshTokenExpDays:   get("REFRESH_TOKEN_EXP_DAYS"),
		},

		Redis: RedisConfig{
			RedisHost:     get("REDIS_HOST"),
			RedisPort:     get("REDIS_PORT"),
			RedisPassword: get("REDIS_PASSWORD"),
		},

		Application: ApplicationConfig{
			Environment: get("ENVIRONMENT"),
		},
	}

}

// Getters

func get(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		log.Panicf("%s is null", envName)
	}
	return env
}
