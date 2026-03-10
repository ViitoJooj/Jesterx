package middleware

import (
	"net/http"
	"strings"

	"github.com/ViitoJooj/Jesterx/internal/config"
)

func allowedOrigin(origin string) bool {
	if origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173" {
		return true
	}
	frontendURL := config.FrontendURL
	if frontendURL == "" {
		return false
	}
	for _, allowed := range strings.Split(frontendURL, ",") {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}
	return false
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Website-Id")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		isPublicRender := strings.HasPrefix(r.URL.Path, "/p/")
		if !isPublicRender {
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; connect-src 'self'")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
