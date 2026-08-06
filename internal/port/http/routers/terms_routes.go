package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterTermsRoutes(mux *http.ServeMux, controller *controllers.TermsController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /terms", wrapHandler(controller.Create, middlewares...))
	mux.Handle("GET /terms", wrapHandler(controller.GetAll, middlewares...))
	mux.Handle("GET /terms/uuid", wrapHandler(controller.GetByUUID, middlewares...))
}
