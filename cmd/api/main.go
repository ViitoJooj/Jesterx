package main

import (
	"os"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	"github.com/ViitoJooj/Jesterx/pkg/logger"
	postgresql "github.com/ViitoJooj/Jesterx/pkg/postgreSQL"
)

func main() {
	config, err := dotenv.Conn()
	if err != nil {
		logger.Warn(err).Print()
		os.Exit(0)
	}

	db, err := postgresql.Conn(config.PostgreSQL)
	if err != nil {
		logger.Warn(err).Print()
		os.Exit(0)
	}

	err = postgresql.Migrations(db)
	if err != nil {
		logger.Warn(err).Print()
		os.Exit(0)
	}
}
