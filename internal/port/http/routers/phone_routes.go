package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterPhoneRoutes(mux *http.ServeMux, controller *controllers.PhoneController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /phones", wrapHandler(controller.Create, middlewares...))
	mux.Handle("GET /phones", wrapHandler(controller.GetAll, middlewares...))
	mux.Handle("GET /phones/uuid", wrapHandler(controller.GetByUUID, middlewares...))
}
