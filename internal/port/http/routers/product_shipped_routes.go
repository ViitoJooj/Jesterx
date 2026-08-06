package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterProductShippedRoutes(mux *http.ServeMux, controller *controllers.ProductShippedController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /products-shipped", wrapHandler(controller.Create, middlewares...))
}
