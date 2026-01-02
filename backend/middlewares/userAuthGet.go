package middlewares

import (
	"jesterx-core/config"
	"jesterx-core/helpers"
	"jesterx-core/responses"
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
		var dbRole string
		var banned bool
		err = config.DB.QueryRow(`SELECT COALESCE(plan, 'free'), COALESCE(role, 'platform_user'), COALESCE(banned, FALSE) FROM users WHERE id = $1`, userID).Scan(&currentPlan, &dbRole, &banned)
		if err != nil {
			currentPlan = "free"
			dbRole = "platform_user"
		}

		if banned {
			c.AbortWithStatusJSON(http.StatusForbidden, responses.ErrorResponse{
				Success: false,
				Message: "User is banned",
			})
			return
		}

		email := claims["email"].(string)
		role := helpers.ResolvePlatformRole(email, dbRole)
		if role != dbRole {
			_, _ = config.DB.Exec(`UPDATE users SET role = $1 WHERE id = $2`, role, userID)
		}

		user := helpers.UserData{
			Id:          userID,
			Profile_img: claims["profile_img"].(string),
			First_name:  claims["first_name"].(string),
			Last_name:   claims["last_name"].(string),
			Email:       email,
			Role:        role,
			Plan:        currentPlan,
		}

		c.Set("user", user)

		c.Next()
	}
}
