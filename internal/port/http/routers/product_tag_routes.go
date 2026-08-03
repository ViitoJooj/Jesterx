package routers

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/port/http/controllers"
)

func RegisterProductTagRoutes(mux *http.ServeMux, c *controllers.ProductTagController) {
	mux.HandleFunc("POST /product-tags", c.Create)
}
