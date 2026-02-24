package middleware

import (
	"context"
	"net/http"

	"github.com/ViitoJooj/Jesterx/internal/security"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

type contextKey string

const UserIDKey contextKey = "userID"

func UserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(UserIDKey).(string)
	return id, ok
}

func IdentityMiddleware(auth *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			websiteId := r.URL.Query().Get("x-website-id")
			if websiteId == "" {
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := r.Cookie(security.AccessCookieName(websiteId))
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := security.ParseAccessToken(cookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user, err := auth.GetUserByID(claims.Sub)
			if err != nil || user == nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, user.Id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := UserID(r.Context()); !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
