package services

import (
	"jesterx-core/config"
	"jesterx-core/helpers"
	"jesterx-core/models"
	"jesterx-core/responses"
	"jesterx-core/templates"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePageService(c *gin.Context) {
	user := c.MustGet("user").(helpers.UserData)
	tenantID := c.GetString("tenantID")

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

	if len(body.Components) > componentLimit {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Limite de componentes excedido."})
		return
	}

	svelteContent := templates.GetTemplateByType(body.PageType)
	if body.Template != "" {
		svelteContent = body.Template
	}

	if tenantID == "" {
		if fallbackTenant, err := resolveUserTenant(user.Id); err == nil && fallbackTenant != "" {
			tenantID = fallbackTenant
		}
	}

	hasAccess := true
	if tenantID == "" {
		hasAccess = false
	}
	if tenantID != "" {
		ok, _, err := helpers.UserHasTenantAccess(user.Id, tenantID)
		if err != nil {
			c.JSON(500, responses.ErrorResponse{Success: false, Message: "Failed to check permissions."})
			return
		}
		hasAccess = ok
	}
	if !hasAccess {
		c.JSON(403, responses.ErrorResponse{Success: false, Message: "Você precisa criar ou selecionar um site para adicionar páginas."})
		return
	}

	db := config.DB
	pageUUID := uuid.New().String()

	slug := body.PageID
	if slug == "" {
		slug = pageUUID
	}

	plan, err := GetPlanConfig(c.Request.Context(), strings.TrimSpace(user.Plan))
	if err != nil {
		plan = PlanConfig{RouteLimit: 0}
	}

	if plan.RouteLimit > 0 {
		var pageCount int
		if err := db.QueryRow(`SELECT COUNT(*) FROM pages WHERE tenant_id = $1`, tenantID).Scan(&pageCount); err == nil {
			if pageCount >= plan.RouteLimit {
				c.JSON(403, responses.ErrorResponse{Success: false, Message: "Limite de rotas do seu plano atingido."})
				return
			}
		}
	}

	_, err = db.Exec(`INSERT INTO pages (id, tenant_id, name, page_id) VALUES ($1, $2, $3, $4)`, pageUUID, tenantID, body.Name, slug)
	if err != nil {
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Failed to save page in Postgres."})
		return
	}

	now := time.Now().UTC()
	showHeader := true
	if body.ShowHeader != nil {
		showHeader = *body.ShowHeader
	}
	showFooter := true
	if body.ShowFooter != nil {
		showFooter = *body.ShowFooter
	}

	doc := models.PageSvelte{
		ID:         pageUUID,
		TenantID:   tenantID,
		PageID:     slug,
		Svelte:     svelteContent,
		Header:     body.Header,
		Footer:     body.Footer,
		ShowHeader: showHeader,
		ShowFooter: showFooter,
		Components: body.Components,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	_, err = config.MongoClient.Database("genyou").Collection("page_sveltes").InsertOne(c.Request.Context(), doc)

	if err != nil {
		_, _ = db.Exec(`DELETE FROM pages WHERE id = $1`, pageUUID)
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Failed to save page content in MongoDB."})
		return
	}

	_, _ = config.MongoClient.Database("genyou").Collection("theme_store_entries").InsertOne(c.Request.Context(), models.ThemeStoreEntry{
		ID:        pageUUID,
		TenantID:  tenantID,
		PageID:    slug,
		Name:      body.Name,
		Domain:    body.Domain,
		ForSale:   false,
		CreatedAt: now,
		UpdatedAt: now,
	})

	c.JSON(201, responses.CreatePageDTO{
		Id:        pageUUID,
		Page_id:   slug,
		Name:      body.Name,
		Tenant_id: tenantID,
	})
}
