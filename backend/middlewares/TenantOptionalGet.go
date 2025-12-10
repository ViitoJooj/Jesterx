package middlewares

import (
	"database/sql"
	"gen-you-ecommerce/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OptionalTenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantPageID := c.GetHeader("X-Tenant-Page-Id")
		if tenantPageID == "" {
			c.Next()
			return
		}

		var tenantID string
		err := config.DB.QueryRow(`SELECT id FROM tenants WHERE page_id = $1`, tenantPageID).Scan(&tenantID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"success": false, "error": "Tenant not found"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Database error"})
			return
		}

		c.Set("tenantID", tenantID)
		c.Next()
	}
}
