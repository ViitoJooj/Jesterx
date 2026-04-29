package dotenv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var SupabaseURI string
var ResendKey string
var SecretKey string
var SecretKeyExpTime string
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
	SecretKeyExpTime = get("SECRET_KEY_EXP_TIME")
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
