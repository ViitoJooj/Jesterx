package http

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/http/handlers"
)

func NewRouter(authHandler *handlers.AuthHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("GET /api/v1/auth/refresh", authHandler.Refresh)
	return mux
}
