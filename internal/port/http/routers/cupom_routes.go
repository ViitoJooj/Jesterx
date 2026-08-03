package routers

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/port/http/controllers"
)

func RegisterCupomRoutes(mux *http.ServeMux, c *controllers.CupomController) {
	mux.HandleFunc("POST /cupoms", c.Create)
}
