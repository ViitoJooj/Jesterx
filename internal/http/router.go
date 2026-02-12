package http

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/http/handlers"
)

func NewRouter(authHandler *handlers.AuthHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("/api/v1/auth//login", authHandler.Login)
	return mux
}
