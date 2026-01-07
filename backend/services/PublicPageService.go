package services

import (
	"database/sql"
	"jesterx-core/config"
	"jesterx-core/models"
	"jesterx-core/responses"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PublicPageService(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(string)
	pageSlug := c.Param("page_id")

	if pageSlug == "" {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Invalid page id."})
		return
	}

	var page models.Page
	err := config.DB.QueryRow(`
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

	products := []models.Product{}
	cursor, _ := config.MongoClient.Database("genyou").Collection("products").Find(c.Request.Context(), bson.M{
		"tenant_id": tenantID,
		"page_id":   pageSlug,
		"visible":   true,
	})
	defer cursor.Close(c.Request.Context())
	for cursor.Next(c.Request.Context()) {
		var p models.Product
		if err := cursor.Decode(&p); err == nil {
			products = append(products, p)
		}
	}

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"meta": responses.PageDTO{
				Id:         page.Id,
				Tenant_id:  page.TenantId,
				Name:       page.Name,
				Page_id:    page.PageId,
				Domain:     page.Domain.String,
				Theme_id:   page.ThemeId.String,
				Created_at: page.CreatedAt,
				Updated_at: page.UpdatedAt,
			},
			"content": gin.H{
				"id":          doc.ID,
				"svelte":      doc.Svelte,
				"header":      doc.Header,
				"footer":      doc.Footer,
				"show_header": doc.ShowHeader,
				"show_footer": doc.ShowFooter,
				"components":  doc.Components,
			},
			"products": products,
		},
	})
}
