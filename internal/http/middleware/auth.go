package middleware

import (
	"net/http"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	"github.com/ViitoJooj/Jesterx/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func Auth(cfg dotenv.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil || accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		claims, err := jwt.ValidateAccessToken(accessToken, cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("user_id", claims.Subject)
		c.Next()
	}
}
