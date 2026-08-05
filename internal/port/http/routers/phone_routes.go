package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterPhoneRoutes(mux *http.ServeMux, c *controllers.PhoneController) {
	mux.HandleFunc("POST /phones", c.Create)
	mux.HandleFunc("GET /phones", c.GetAll)
	mux.HandleFunc("GET /phones/uuid", c.GetByUUID)
}
