package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterProductRoutes(mux *http.ServeMux, c *controllers.ProductController) {
	mux.HandleFunc("POST /products", c.Create)
}
