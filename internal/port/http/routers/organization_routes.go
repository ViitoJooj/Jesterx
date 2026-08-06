package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterOrganizationRoutes(mux *http.ServeMux, controller *controllers.OrganizationController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /organizations", wrapHandler(controller.Create, middlewares...))
	mux.Handle("GET /organizations", wrapHandler(controller.GetAll, middlewares...))
	mux.Handle("GET /organizations/uuid", wrapHandler(controller.GetByUUID, middlewares...))
}
