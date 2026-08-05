package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterTermsRoutes(mux *http.ServeMux, c *controllers.TermsController) {
	mux.HandleFunc("POST /terms", c.Create)
	mux.HandleFunc("GET /terms", c.GetAll)
	mux.HandleFunc("GET /terms/uuid", c.GetByUUID)
}
