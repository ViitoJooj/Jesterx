package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterProductShippedRoutes(mux *http.ServeMux, c *controllers.ProductShippedController) {
	mux.HandleFunc("POST /products-shipped", c.Create)
}
