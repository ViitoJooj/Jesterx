package services

import (
	"database/sql"
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/helpers"
	"gen-you-ecommerce/models"
	"gen-you-ecommerce/responses"
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

	if len(body.Svelte) > 500_000 {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Svelte content is too large."})
		return
	}

	hasAccess, _, err := userHasTenantAccess(user.Id, tenantID)
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
		Svelte:    body.Svelte,
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

func userHasTenantAccess(userID, tenantID string) (bool, string, error) {
	db := config.DB
	var role string

	err := db.QueryRow(`SELECT role FROM tenant_users WHERE user_id = $1 AND tenant_id = $2`, userID, tenantID).Scan(&role)

	if err == sql.ErrNoRows {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}

	return true, role, nil
}
