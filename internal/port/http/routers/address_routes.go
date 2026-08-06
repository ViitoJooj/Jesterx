package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterAddressRoutes(mux *http.ServeMux, controller *controllers.AddressController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /addresses", wrapHandler(controller.Create, middlewares...))
	mux.Handle("GET /addresses", wrapHandler(controller.GetAll, middlewares...))
	mux.Handle("GET /addresses/uuid", wrapHandler(controller.GetByUUID, middlewares...))
}
