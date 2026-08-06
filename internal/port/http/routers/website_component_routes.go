package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterWebsiteComponentRoutes(mux *http.ServeMux, controller *controllers.WebsiteComponentController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /website-components", wrapHandler(controller.Create, middlewares...))
}
