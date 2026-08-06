package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterRbacRoutes(mux *http.ServeMux, controller *controllers.RbacController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /rbacs", wrapHandler(controller.Create, middlewares...))
	mux.Handle("GET /rbacs", wrapHandler(controller.GetAll, middlewares...))
	mux.Handle("GET /rbacs/uuid", wrapHandler(controller.GetByUUID, middlewares...))
}
