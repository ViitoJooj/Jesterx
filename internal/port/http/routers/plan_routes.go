package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterPlanRoutes(mux *http.ServeMux, c *controllers.PlanController) {
	mux.HandleFunc("POST /plans", c.Create)
}
