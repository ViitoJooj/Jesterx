package dotenv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Supabase_uri string

func Set() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Panic("Error on get .env")
	}

	Supabase_uri = get("Supabase_uri")

}

func get(env_name string) string {

	env := os.Getenv(env_name)
	if env == "" {
		log.Panicf("%s is null", env)
	}

	return env
}
