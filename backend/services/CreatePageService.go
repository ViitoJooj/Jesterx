package services

import (
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/helpers"
	"gen-you-ecommerce/models"
	"gen-you-ecommerce/responses"
	"gen-you-ecommerce/templates"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePageService(c *gin.Context) {
	user := c.MustGet("user").(helpers.UserData)
	tenantID := c.MustGet("tenantID").(string)

	var body models.CreatePageModels

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Invalid request body."})
		return
	}

	validTypes := map[string]bool{"landing": true, "ecommerce": true, "software": true, "video": true}
	if !validTypes[body.PageType] {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Invalid page type."})
		return
	}

	svelteContent := templates.GetTemplateByType(body.PageType)
	if body.Template != "" {
		svelteContent = body.Template
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
	pageUUID := uuid.New().String()

	slug := body.PageID
	if slug == "" {
		slug = pageUUID
	}

	_, err = db.Exec(`INSERT INTO pages (id, tenant_id, name, page_id) VALUES ($1, $2, $3, $4)`, pageUUID, tenantID, body.Name, slug)
	if err != nil {
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Failed to save page in Postgres."})
		return
	}

	now := time.Now().UTC()
	doc := models.PageSvelte{
		ID:        pageUUID,
		TenantID:  tenantID,
		PageID:    slug,
		Svelte:    svelteContent,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = config.MongoClient.Database("genyou").Collection("page_sveltes").InsertOne(c.Request.Context(), doc)

	if err != nil {
		_, _ = db.Exec(`DELETE FROM pages WHERE id = $1`, pageUUID)
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Failed to save page content in MongoDB."})
		return
	}

	c.JSON(201, responses.CreatePageDTO{
		Id:        pageUUID,
		Page_id:   slug,
		Name:      body.Name,
		Tenant_id: tenantID,
	})
}
