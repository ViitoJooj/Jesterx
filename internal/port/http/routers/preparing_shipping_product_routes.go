package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterPreparingShippingProductRoutes(mux *http.ServeMux, c *controllers.PreparingShippingProductController) {
	mux.HandleFunc("POST /preparing-shipping-products", c.Create)
}
