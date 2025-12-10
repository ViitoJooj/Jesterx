package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JwtSecret string
var PostgresUri string
var MongoUri string
var Gin_mode string

func Load() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	postgres_user := os.Getenv("POSTGRES_USER")
	if postgres_user == "" {
		log.Fatal("Error on get POSTGRES_USER")
	}
	postgres_password := os.Getenv("POSTGRES_PASSWORD")
	if postgres_password == "" {
		log.Fatal("Error on get POSTGRES_PASSWORD")
	}
	postgres_db := os.Getenv("POSTGRES_DB")
	if postgres_db == "" {
		log.Fatal("Error on get POSTGRES_DB")
	}
	postgres_port := os.Getenv("POSTGRES_PORT")
	if postgres_port == "" {
		log.Fatal("Error on get POSTGRES_PORT")
	}
	postgres_host := os.Getenv("POSTGRES_HOST")
	if postgres_host == "" {
		log.Fatal("Error on get POSTGRES_HOST")
	}
	postgres_ssl := os.Getenv("POSTGRES_SSL")
	if postgres_ssl == "" {
		log.Fatal("Error on get POSTGRES_SSL")
	}

	PostgresUri = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", postgres_user, postgres_password, postgres_host, postgres_port, postgres_db, postgres_ssl)

	mongo_user := os.Getenv("MONGO_USER")
	if mongo_user == "" {
		log.Fatal("Error on get MONGO_USER")
	}
	mongo_password := os.Getenv("MONGO_PASSWORD")
	if mongo_password == "" {
		log.Fatal("Error on get MONGO_PASSWORD")
	}
	mongo_port := os.Getenv("MONGO_PORT")
	if mongo_port == "" {
		log.Fatal("Error on get MONGO_PORT")
	}
	mongo_host := os.Getenv("MONGO_HOST")
	if mongo_host == "" {
		log.Fatal("Error on get MONGO_HOST")
	}

	MongoUri = fmt.Sprintf("mongodb://%s:%s@%s:%s", mongo_user, mongo_password, mongo_host, mongo_port)

	JwtSecret = os.Getenv("JWT_SECRET")
	if JwtSecret == "" {
		log.Fatal("Error on get JWT_SECRET")
	}

	Gin_mode = os.Getenv("GIN_MODE")
	if Gin_mode == "" {
		log.Fatal("Error on get GIN_MODE")
	}
}
