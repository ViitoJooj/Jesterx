package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterAuthRoutes(mux *http.ServeMux, controller *controllers.AuthController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /auth/register", wrapHandler(controller.Register, middlewares...))
	mux.Handle("POST /auth/login", wrapHandler(controller.Login, middlewares...))
}
