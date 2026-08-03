package routers

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/port/http/controllers"
)

func RegisterRbacRoutes(mux *http.ServeMux, c *controllers.RbacController) {
	mux.HandleFunc("POST /rbacs", c.Create)
	mux.HandleFunc("GET /rbacs", c.GetAll)
	mux.HandleFunc("GET /rbacs/uuid", c.GetByUUID)
}
