package dotenv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var SupabaseURI string
var ResendKey string
var SecretKey string
var AccessTokenExpMinutes string
var RefreshTokenExpDays string
var RedisHost string
var RedisPort string
var RedisPassword string

func Conn() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Error loading .env: %v", err)
	}

	SupabaseURI = get("SUPABASE_URI")
	ResendKey = get("RESEND_API_KEY")
	SecretKey = get("SECRET_KEY")
	AccessTokenExpMinutes = get("ACCESS_TOKEN_EXP_MINUTES")
	RefreshTokenExpDays = get("REFRESH_TOKEN_EXP_DAYS")
	RedisHost = get("REDIS_HOST")
	RedisPort = get("REDIS_PORT")
	RedisPassword = get("REDIS_PASSWORD")
}

// Getters

func get(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		log.Panicf("%s is null", envName)
	}
	return env
}
