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

func LoadEnv() {
	_ = godotenv.Load("../.env")

	PGCNN.User = mustGetenv("POSTGRES_USER")
	PGCNN.Password = mustGetenv("POSTGRES_PASSWORD")
	PGCNN.DBName = mustGetenv("POSTGRES_DB")
	PGCNN.Port = mustGetenv("POSTGRES_PORT")
	PGCNN.Host = mustGetenv("POSTGRES_HOST")
	PGCNN.SSLMode = mustGetenv("POSTGRES_SSL")
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatal("Error on get " + key)
	}
	return v
}
