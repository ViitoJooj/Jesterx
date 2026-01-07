package services

import (
	"database/sql"
	"jesterx-core/config"
	"jesterx-core/helpers"
	"jesterx-core/models"
	"jesterx-core/responses"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetPageService(c *gin.Context) {
	user := c.MustGet("user").(helpers.UserData)
	tenantID := c.MustGet("tenantID").(string)
	pageSlug := c.Param("page_id")

	if pageSlug == "" {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Invalid page id."})
		return
	}

	hasAccess, _, err := helpers.UserHasTenantAccess(user.Id, tenantID)
	if err != nil {
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Failed to check permissions."})
		return
	}
	if !hasAccess {
		c.JSON(403, responses.ErrorResponse{Success: false, Message: "You do not belong to this site (tenant)."})
		return
	}

	db := config.DB

	var page models.Page
	err = db.QueryRow(`
        SELECT id, tenant_id, name, page_id, domain, theme_id, created_at, updated_at
        FROM pages
        WHERE tenant_id = $1 AND page_id = $2
    `, tenantID, pageSlug).Scan(
		&page.Id,
		&page.TenantId,
		&page.Name,
		&page.PageId,
		&page.Domain,
		&page.ThemeId,
		&page.CreatedAt,
		&page.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(404, responses.ErrorResponse{Success: false, Message: "Page not found."})
		return
	}
	if err != nil {
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Failed to load page."})
		return
	}

	c.JSON(200, responses.PageDTO{
		Id:         page.Id,
		Tenant_id:  page.TenantId,
		Name:       page.Name,
		Page_id:    page.PageId,
		Domain:     page.Domain.String,
		Theme_id:   page.ThemeId.String,
		Created_at: page.CreatedAt,
		Updated_at: page.UpdatedAt,
	})
}

func GetRawSveltePageService(c *gin.Context) {
	user := c.MustGet("user").(helpers.UserData)
	tenantID := c.MustGet("tenantID").(string)
	pageSlug := c.Param("page_id")

	if pageSlug == "" {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Invalid page id."})
		return
	}

	hasAccess, _, err := helpers.UserHasTenantAccess(user.Id, tenantID)
	if err != nil {
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Failed to check permissions."})
		return
	}
	if !hasAccess {
		c.JSON(403, responses.ErrorResponse{Success: false, Message: "You do not belong to this site (tenant)."})
		return
	}

	coll := config.MongoClient.Database("genyou").Collection("page_sveltes")

	var doc models.PageSvelte
	err = coll.FindOne(c.Request.Context(), bson.M{
		"tenant_id": tenantID,
		"page_id":   pageSlug,
	}).Decode(&doc)

	if err == mongo.ErrNoDocuments {
		c.JSON(404, responses.ErrorResponse{Success: false, Message: "Page content not found."})
		return
	}
	if err != nil {
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Failed to load page content."})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"id":          doc.ID,
			"tenant_id":   doc.TenantID,
			"page_id":     doc.PageID,
			"svelte":      doc.Svelte,
			"header":      doc.Header,
			"footer":      doc.Footer,
			"show_header": doc.ShowHeader,
			"show_footer": doc.ShowFooter,
			"components":  doc.Components,
			"created_at":  doc.CreatedAt,
			"updated_at":  doc.UpdatedAt,
		},
	})
}
