package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterProductTagRoutes(mux *http.ServeMux, controller *controllers.ProductTagController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /product-tags", wrapHandler(controller.Create, middlewares...))
}
