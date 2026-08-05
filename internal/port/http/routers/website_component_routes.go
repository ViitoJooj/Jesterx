package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterWebsiteComponentRoutes(mux *http.ServeMux, c *controllers.WebsiteComponentController) {
	mux.HandleFunc("POST /website-components", c.Create)
}
