package services

import (
	"jesterx-core/config"
	"jesterx-core/helpers"
	"jesterx-core/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PageID    string `json:"page_id"`
	Domain    string `json:"domain"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ListPagesService(c *gin.Context) {
	user := c.MustGet("user").(helpers.UserData)
	tenantID := c.GetString("tenantID")

	if tenantID == "" {
		if fallbackTenant, err := resolveUserTenant(user.Id); err == nil && fallbackTenant != "" {
			tenantID = fallbackTenant
		}
	}

	hasAccess := true
	if tenantID == "" {
		hasAccess = false
	} else {
		var err error
		hasAccess, _, err = helpers.UserHasTenantAccess(user.Id, tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to check permissions."})
			return
		}
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, responses.ErrorResponse{Success: false, Message: "Você precisa criar ou selecionar um site para listar páginas."})
		return
	}

	db := config.DB
	rows, err := db.Query(`
		SELECT id, name, page_id, COALESCE(domain, ''), created_at, updated_at
		FROM pages
		WHERE tenant_id = $1
		ORDER BY created_at DESC
	`, tenantID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to fetch pages."})
		return
	}
	defer rows.Close()

	pages := []PageItem{}
	for rows.Next() {
		var page PageItem
		err := rows.Scan(&page.ID, &page.Name, &page.PageID, &page.Domain, &page.CreatedAt, &page.UpdatedAt)
		if err != nil {
			continue
		}
		pages = append(pages, page)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pages,
	})
}
