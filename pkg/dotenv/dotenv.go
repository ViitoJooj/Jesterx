package dotenv

import (
	"os"

	"github.com/joho/godotenv"
)

func Conn() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		Application: Application{
			Port:       os.Getenv("PORT"),
			Enviroment: os.Getenv("ENVIROMENT"),
			ViewUrl:    os.Getenv("VIEW_URL"),
			DaemonUrl:  os.Getenv("DAEMON_URL"),
		},
		PostgreSQL: PostgreSQL{
			URI:      os.Getenv("POSTGRES_URI"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
		},
		Security: Security{
			PasetoSecretKey: os.Getenv("PASETO_SECRET_KEY"),
		},
	}, nil
}
