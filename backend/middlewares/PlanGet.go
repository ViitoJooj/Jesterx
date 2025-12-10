package middlewares

import (
	"database/sql"
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PlanLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(helpers.UserData)
		tenantID := c.MustGet("tenantID").(string)

		maxPages := maxPagesForPlan(user.Plan)
		if maxPages == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Your current plan does not allow creating pages.",
			})
			return
		}

		var count int
		err := config.DB.QueryRow(
			`SELECT COUNT(*) FROM pages WHERE tenant_id = $1`,
			tenantID,
		).Scan(&count)

		if err != nil && err != sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check page limit.",
			})
			return
		}

		if count >= maxPages {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Page limit reached for your current subscription plan.",
			})
			return
		}

		c.Next()
	}
}

func maxPagesForPlan(plan string) int {
	switch plan {
	case "free":
		return 0
	case "business":
		return 1
	case "pro":
		return 10
	case "enterprise":
		return 50
	default:
		return 0
	}
}
