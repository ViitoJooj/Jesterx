package middlewares

import (
	"database/sql"
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

var planLimits = map[string]int{
	"free":       0,
	"business":   1,
	"pro":        10,
	"enterprise": 50,
}

func PlanMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAny, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "User not in context",
			})
			return
		}

		user := userAny.(helpers.UserData)

		// Sempre pega o plano do banco pra garantir que está atualizado
		var plan string
		err := config.DB.QueryRow(`SELECT plan FROM users WHERE id = $1`, user.Id).Scan(&plan)
		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "User not found",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Database error while reading user plan",
			})
			return
		}

		maxSites, ok := planLimits[plan]
		if !ok {
			maxSites = 0
		}

		var currentSites int
		err = config.DB.QueryRow(`
			SELECT COUNT(*) 
			FROM tenant_users 
			WHERE user_id = $1 AND role = 'owner'
		`, user.Id).Scan(&currentSites)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Database error while counting user sites",
			})
			return
		}

		if currentSites >= maxSites {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Plan limit reached for creating new sites",
				"plan":    plan,
			})
			return
		}

		c.Next()
	}
}
