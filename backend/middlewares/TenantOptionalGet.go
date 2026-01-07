package middlewares

import (
	"database/sql"
	"jesterx-core/config"
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
			// Se o tenant não existir, apenas ignore o header para não quebrar rotas públicas/sem tenant
			if err != sql.ErrNoRows {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Database error"})
				return
			}
			c.Next()
			return
		}

		c.Set("tenantID", tenantID)
		c.Next()
	}
}
