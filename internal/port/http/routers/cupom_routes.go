package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterCupomRoutes(mux *http.ServeMux, controller *controllers.CupomController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /cupoms", wrapHandler(controller.Create, middlewares...))
}
