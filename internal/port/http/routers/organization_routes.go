package routers

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/port/http/controllers"
)

func RegisterOrganizationRoutes(mux *http.ServeMux, c *controllers.OrganizationController) {
	mux.HandleFunc("POST /organizations", c.Create)
	mux.HandleFunc("GET /organizations", c.GetAll)
	mux.HandleFunc("GET /organizations/uuid", c.GetByUUID)
}
