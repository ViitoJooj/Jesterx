package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterCupomRoutes(mux *http.ServeMux, c *controllers.CupomController) {
	mux.HandleFunc("POST /cupoms", c.Create)
}
