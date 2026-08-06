package routers

import "net/http"

func wrapHandler(handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) http.Handler {
	var h http.Handler = handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
