package main

import (
	"os"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	"github.com/ViitoJooj/Jesterx/pkg/logger"
	"github.com/ViitoJooj/Jesterx/pkg/postgresql"
	"github.com/ViitoJooj/Jesterx/pkg/server"
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

	server.Start(config.Application.Port)
}
