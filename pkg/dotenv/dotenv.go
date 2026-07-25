package dotenv

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Conn() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &Config{
		PostgreSQL: PostgreSQL{
			URI:      os.Getenv("POSTGRES_URI"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
		},
	}, nil
}
