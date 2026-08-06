package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

func RegisterPlanRoutes(mux *http.ServeMux, controller *controllers.PlanController, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("POST /plans", wrapHandler(controller.Create, middlewares...))
}
