package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterStorageProductRoutes(mux *http.ServeMux, c *controllers.StorageProductController) {
	mux.HandleFunc("POST /storage-products", c.Create)
}
