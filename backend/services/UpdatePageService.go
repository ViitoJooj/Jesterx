package services

import (
	"database/sql"
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/helpers"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdatePageService(c *gin.Context) {
	user := c.MustGet("user").(helpers.UserData)
	tenantID := c.MustGet("tenantID").(string)
	currentSlug := c.Param("page_id")

	var body struct {
		Name   *string `json:"name"`
		PageID *string `json:"page_id"`
		Svelte *string `json:"svelte"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body."})
		return
	}

	if body.Name == nil && body.PageID == nil && body.Svelte == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing to update."})
		return
	}

	if body.Svelte != nil && len(*body.Svelte) > 200_000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Svelte content is too large."})
		return
	}

	hasAccess, role, err := userHasTenantAccess(user.Id, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions."})
		return
	}
	if !hasAccess || (role != "owner" && role != "admin" && role != "editor") {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to edit pages."})
		return
	}

	db := config.DB

	var pageID string
	var name string
	var slug string

	err = db.QueryRow(`
        SELECT id, name, page_id
        FROM pages
        WHERE tenant_id = $1 AND page_id = $2
    `, tenantID, currentSlug).Scan(&pageID, &name, &slug)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load page."})
		return
	}

	if body.Name != nil {
		name = *body.Name
	}
	if body.PageID != nil && *body.PageID != "" {
		slug = *body.PageID
	}

	_, err = db.Exec(`
        UPDATE pages
        SET name = $1,
            page_id = $2,
            updated_at = NOW()
        WHERE id = $3 AND tenant_id = $4
    `, name, slug, pageID, tenantID)

	if err != nil {
		if strings.Contains(err.Error(), "pages_unique_tenant_page") {
			c.JSON(http.StatusConflict, gin.H{"error": "A page with this slug already exists in this site."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update page in Postgres."})
		return
	}

	updateDoc := bson.M{
		"page_id":    slug,
		"updated_at": time.Now().UTC(),
	}
	if body.Svelte != nil {
		updateDoc["svelte"] = *body.Svelte
	}

	_, err = config.MongoClient.
		Database("genyou").
		Collection("page_sveltes").
		UpdateOne(
			c.Request.Context(),
			bson.M{"_id": pageID},
			bson.M{"$set": updateDoc},
		)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update page content in MongoDB."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       pageID,
		"page_id":  slug,
		"name":     name,
		"tenantId": tenantID,
	})
}
