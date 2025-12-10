package services

import (
	"database/sql"
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/helpers"
	"gen-you-ecommerce/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePageService(c *gin.Context) {
	user := c.MustGet("user").(helpers.UserData)
	tenantID := c.MustGet("tenantID").(string)

	var body struct {
		Name   string `json:"name" binding:"required"`
		PageID string `json:"page_id"`
		Svelte string `json:"svelte" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body."})
		return
	}

	if len(body.Svelte) > 200_000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Svelte content is too large."})
		return
	}

	hasAccess, _, err := userHasTenantAccess(user.Id, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions."})
		return
	}
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not belong to this site (tenant)."})
		return
	}

	db := config.DB
	pageUUID := uuid.New().String()

	slug := body.PageID
	if slug == "" {
		slug = pageUUID
	}

	_, err = db.Exec(`
        INSERT INTO pages (id, tenant_id, name, page_id)
        VALUES ($1, $2, $3, $4)
    `, pageUUID, tenantID, body.Name, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save page in Postgres.",
			"details": err.Error(),
		})
		return
	}

	now := time.Now().UTC()
	doc := models.PageSvelte{
		ID:        pageUUID,
		TenantID:  tenantID,
		PageID:    slug,
		Svelte:    body.Svelte,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = config.MongoClient.
		Database("genyou").
		Collection("page_sveltes").
		InsertOne(c.Request.Context(), doc)
	if err != nil {
		_, _ = db.Exec(`DELETE FROM pages WHERE id = $1`, pageUUID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save page content in MongoDB."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       pageUUID,
		"page_id":  slug,
		"name":     body.Name,
		"tenantId": tenantID,
	})
}

func userHasTenantAccess(userID, tenantID string) (bool, string, error) {
	db := config.DB
	var role string

	err := db.QueryRow(`
        SELECT role
        FROM tenant_users
        WHERE user_id = $1 AND tenant_id = $2
    `, userID, tenantID).Scan(&role)

	if err == sql.ErrNoRows {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}

	return true, role, nil
}
