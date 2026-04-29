package dotenv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var SupabaseURI string
var ResendKey string
var SecretKey string

func Set() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Error loading .env: %v", err)
	}

	SupabaseURI = get("SUPABASE_URI")
	ResendKey = get("RESEND_KEY")
	SecretKey = get("SECRET_KEY")
}

func get(envName string) string {
	env := os.Getenv(envName)
	if env == "" {
		log.Panicf("%s is null", envName)
	}
	return env
}
