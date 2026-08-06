package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ViitoJooj/verkoupe/pkg/token"
)

type contextKey string

const UserUUIDKey contextKey = "user_uuid"
const WebsiteUUIDKey contextKey = "website_uuid"
const RoleKey contextKey = "role"

func AuthMiddleware(pasetoSecret string) func(http.Handler) http.Handler {
	secret := []byte(pasetoSecret)
	if len(secret) < 32 {
		padded := make([]byte, 32)
		copy(padded, secret)
		secret = padded
	} else if len(secret) > 32 {
		secret = secret[:32]
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path

			// Public routes that don't require authentication
			if isPublicRoute(path, r.Method) {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeUnauthorized(w, "RAX-012", "unauthorized")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == authHeader {
				writeUnauthorized(w, "RAX-012", "unauthorized")
				return
			}

			claims, err := token.Parse(tokenStr, secret)
			if err != nil {
				writeUnauthorized(w, "RBX-009", "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), UserUUIDKey, claims.UserUUID)
			ctx = context.WithValue(ctx, WebsiteUUIDKey, claims.WebSiteUUID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isPublicRoute(path string, method string) bool {
	publicRoutes := map[string]bool{
		"POST /auth/register": true,
		"POST /auth/login":    true,
		"GET /health":         true,
	}

	return publicRoutes[method+" "+path]
}

func writeUnauthorized(w http.ResponseWriter, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"` + code + `","message":"` + message + `"}`))
}

func GetUserUUID(r *http.Request) string {
	if v, ok := r.Context().Value(UserUUIDKey).(string); ok {
		return v
	}
	return ""
}

func GetWebsiteUUID(r *http.Request) string {
	if v, ok := r.Context().Value(WebsiteUUIDKey).(string); ok {
		return v
	}
	return ""
}
