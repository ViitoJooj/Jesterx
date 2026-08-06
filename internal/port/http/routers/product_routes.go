package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterProductRoutes(mux *http.ServeMux, controller *controllers.ProductController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /products", wrapHandler(controller.Create, middlewares...))
}
