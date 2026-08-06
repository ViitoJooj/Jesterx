package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterTermsAcceptedRoutes(mux *http.ServeMux, controller *controllers.TermsAcceptedController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /terms-accepted", wrapHandler(controller.Create, middlewares...))
	mux.Handle("GET /terms-accepted", wrapHandler(controller.GetAll, middlewares...))
	mux.Handle("GET /terms-accepted/uuid", wrapHandler(controller.GetByUUID, middlewares...))
}
