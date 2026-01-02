package middlewares

import (
	"net/http"

	"jesterx-core/helpers"
	"jesterx-core/responses"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAny, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responses.ErrorResponse{
				Success: false,
				Message: "Unauthorized",
			})
			return
		}

		user, ok := userAny.(helpers.UserData)
		if !ok || user.Role != "platform_admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, responses.ErrorResponse{
				Success: false,
				Message: "Admin access required",
			})
			return
		}

		c.Next()
	}
}
