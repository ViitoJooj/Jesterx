package routers

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/port/http/controllers"
)

func RegisterStorageProductRoutes(mux *http.ServeMux, c *controllers.StorageProductController) {
	mux.HandleFunc("POST /storage-products", c.Create)
}
