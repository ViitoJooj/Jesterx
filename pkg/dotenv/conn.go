package dotenv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Postgres
var Postgres_host string
var Postgres_port string
var Postgres_user string
var Postgres_password string
var Postgres_ssr string

// Resend
var Resend_key string

// Stripe
var Stripe_secret string
var Stripe_public string

// Application
var IsDev bool
var Frontend_url string
var Backend_url string

func Conn() {
	godotenv.Load(".env")

	Postgres_host = get("POSTGRES_HOST")
	Postgres_port = get("POSTGRES_PORT")
	Postgres_user = get("POSTGRES_USER")
	Postgres_password = get("POSTGRES_PASSWORD")
	Postgres_ssr = get("POSTGRES_SSR")

}

func getWithDefault(envName, defaultVal string) string {
	v := os.Getenv(envName)
	if v == "" {
		return defaultVal
	}
	return v
}

func get(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		log.Panicf("%s is null", envName)
	}
	return env
}
