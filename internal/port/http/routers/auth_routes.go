package routers

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/port/http/controllers"
)

func RegisterAuthRoutes(mux *http.ServeMux, c *controllers.AuthController) {
	mux.HandleFunc("POST /auth/register", c.Register)
}
