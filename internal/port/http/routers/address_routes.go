package routers

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/port/http/controllers"
)

func RegisterAddressRoutes(mux *http.ServeMux, c *controllers.AddressController) {
	mux.HandleFunc("POST /addresses", c.Create)
	mux.HandleFunc("GET /addresses", c.GetAll)
	mux.HandleFunc("GET /addresses/uuid", c.GetByUUID)
}
