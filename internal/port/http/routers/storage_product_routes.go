package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterStorageProductRoutes(mux *http.ServeMux, controller *controllers.StorageProductController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /storage-products", wrapHandler(controller.Create, middlewares...))
}
