package routers

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/port/http/controllers"
)

func RegisterWebsiteRoutes(mux *http.ServeMux, c *controllers.WebsiteController) {
	mux.HandleFunc("POST /websites", c.Create)
}
