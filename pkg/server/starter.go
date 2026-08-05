package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ViitoJooj/verkoupe/pkg/logger"
)

type Conf struct {
	Handler http.Handler
}

func Start(port string, handler http.Handler) {
	if port == "" {
		err := fmt.Errorf("Application port cannot be null.")
		logger.Warn(err).Print()
		os.Exit(0)
	}

	log.Println("Server started in " + port)
	err := http.ListenAndServe(":"+port, handler)
	logger.Warn(err).Print()
	os.Exit(0)
}
