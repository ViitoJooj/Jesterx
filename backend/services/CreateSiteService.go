package services

import (
	"database/sql"
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/helpers"
	"gen-you-ecommerce/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateSiteService(c *gin.Context) {
	user := c.MustGet("user").(helpers.UserData)

	var body struct {
		Name string `json:"name" binding:"required"`
		Slug string `json:"slug" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Invalid request body."})
		return
	}

	maxSites := maxSitesForPlan(user.Plan)
	if maxSites == 0 {
		c.JSON(http.StatusForbidden, responses.ErrorResponse{Success: false, Message: "Your plan does not allow creating sites."})
		return
	}

	db := config.DB

	var currentCount int
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM tenant_users tu
        JOIN tenants t ON t.id = tu.tenant_id
        WHERE tu.user_id = $1 AND tu.role = 'owner'
    `, user.Id).Scan(&currentCount)

	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to check sites limit."})
		return
	}

	if currentCount >= maxSites {
		c.JSON(http.StatusForbidden, responses.ErrorResponse{Success: false, Message: "Site limit reached for your current plan."})
		return
	}

	tenantID := uuid.New().String()

	_, err = db.Exec(`INSERT INTO tenants (id, name, page_id) VALUES ($1, $2, $3)`, tenantID, body.Name, body.Slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to create site (tenant)."})
		return
	}

	_, err = db.Exec(`INSERT INTO tenant_users (tenant_id, user_id, role) VALUES ($1, $2, 'owner')`, tenantID, user.Id)
	if err != nil {
		_, _ = db.Exec(`DELETE FROM tenants WHERE id = $1`, tenantID)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to link user to site."})
		return
	}

	c.JSON(201, gin.H{
		"id":      tenantID,
		"name":    body.Name,
		"slug":    body.Slug,
		"message": "Site created successfully.",
	})
}

func maxSitesForPlan(plan string) int {
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
