package main

import (
	"net/http"
	"os"

	"github.com/ViitoJooj/verkoupe/internal/port/http/routers"
	"github.com/ViitoJooj/verkoupe/pkg/dotenv"
	"github.com/ViitoJooj/verkoupe/pkg/logger"
	"github.com/ViitoJooj/verkoupe/pkg/postgresql"
	"github.com/ViitoJooj/verkoupe/pkg/server"
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

	mux := http.NewServeMux()
	routers.Register(mux, routers.NewControllers(db))

	server.Start(config.Application.Port, mux)
}
