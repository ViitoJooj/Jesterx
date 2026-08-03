package routers

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/port/http/controllers"
)

func RegisterTermsAcceptedRoutes(mux *http.ServeMux, c *controllers.TermsAcceptedController) {
	mux.HandleFunc("POST /terms-accepted", c.Create)
	mux.HandleFunc("GET /terms-accepted", c.GetAll)
	mux.HandleFunc("GET /terms-accepted/uuid", c.GetByUUID)
}
