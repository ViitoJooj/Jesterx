package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterWebsiteRoutes(mux *http.ServeMux, c *controllers.WebsiteController) {
	mux.HandleFunc("POST /websites", c.Create)
}
