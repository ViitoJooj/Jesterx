package routers

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/port/http/controllers"
)

func RegisterProductRoutes(mux *http.ServeMux, c *controllers.ProductController) {
	mux.HandleFunc("POST /products", c.Create)
}
