package middlewares

import (
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/helpers"
	"gen-you-ecommerce/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("auth_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responses.ErrorResponse{
				Success: false,
				Message: "Missing auth token",
			})
			return
		}

		claims, err := helpers.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responses.ErrorResponse{
				Success: false,
				Message: "Invalid or expired auth token",
			})
			return
		}

		userID := claims["sub"].(string)

		var currentPlan string
		err = config.DB.QueryRow(`SELECT COALESCE(plan, 'free') FROM users WHERE id = $1`, userID).Scan(&currentPlan)
		if err != nil {
			currentPlan = "free"
		}

		user := helpers.UserData{
			Id:          userID,
			Profile_img: claims["profile_img"].(string),
			First_name:  claims["first_name"].(string),
			Last_name:   claims["last_name"].(string),
			Email:       claims["email"].(string),
			Role:        claims["role"].(string),
			Plan:        currentPlan,
		}

		c.Set("user", user)

		c.Next()
	}
}
