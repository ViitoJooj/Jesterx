package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterPreparingShippingProductRoutes(mux *http.ServeMux, controller *controllers.PreparingShippingProductController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /preparing-shipping-products", wrapHandler(controller.Create, middlewares...))
}
