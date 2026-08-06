package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterWebsiteRoutes(mux *http.ServeMux, controller *controllers.WebsiteController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /websites", wrapHandler(controller.Create, middlewares...))
	mux.Handle("GET /websites/{uuid}", wrapHandler(controller.GetByUUID, middlewares...))
	mux.Handle("GET /websites", wrapHandler(controller.ListAll, middlewares...))
	mux.Handle("GET /websites/owner", wrapHandler(controller.ListByOwner, middlewares...))
	mux.Handle("PUT /websites/{uuid}", wrapHandler(controller.Update, middlewares...))
	mux.Handle("DELETE /websites/{uuid}", wrapHandler(controller.Delete, middlewares...))
}
