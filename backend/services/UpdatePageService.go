package services

import (
	"database/sql"
	"jesterx-core/config"
	"jesterx-core/helpers"
	"jesterx-core/responses"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdatePageService(c *gin.Context) {
	user := c.MustGet("user").(helpers.UserData)
	tenantID := c.MustGet("tenantID").(string)
	currentSlug := c.Param("page_id")

	var body struct {
		Name       *string  `json:"name"`
		PageID     *string  `json:"page_id"`
		Svelte     *string  `json:"svelte"`
		Header     *string  `json:"header"`
		Footer     *string  `json:"footer"`
		Components []string `json:"components"`
		ShowHeader *bool    `json:"show_header"`
		ShowFooter *bool    `json:"show_footer"`
		Domain     *string  `json:"domain"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Invalid request body."})
		return
	}

	if body.Name == nil && body.PageID == nil && body.Svelte == nil && body.Header == nil && body.Footer == nil && body.Domain == nil && body.ShowFooter == nil && body.ShowHeader == nil && body.Components == nil {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Nothing to update."})
		return
	}

	if body.Svelte != nil && len(*body.Svelte) > 200_000 {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Svelte content is too large."})
		return
	}

	if body.Components != nil && len(body.Components) > componentLimit {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Limite de componentes excedido."})
		return
	}

	hasAccess, role, err := helpers.UserHasTenantAccess(user.Id, tenantID)
	if err != nil {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Failed to check permissions."})
		return
	}
	if !hasAccess || (role != "owner" && role != "admin" && role != "editor") {
		c.JSON(403, responses.ErrorResponse{Success: false, Message: "You do not have permission to edit pages."})
		return
	}

	db := config.DB

	var pageID string
	var name string
	var slug string
	var domain sql.NullString

	err = db.QueryRow(`SELECT id, name, page_id, domain FROM pages WHERE tenant_id = $1 AND page_id = $2`, tenantID, currentSlug).Scan(&pageID, &name, &slug, &domain)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, responses.ErrorResponse{Success: false, Message: "Page not found."})
			return
		}
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Failed to load page."})
		return
	}

	if body.Name != nil {
		name = *body.Name
	}
	if body.PageID != nil && *body.PageID != "" {
		slug = *body.PageID
	}
	if body.Domain != nil {
		d := strings.TrimSpace(*body.Domain)
		if d == "" {
			domain = sql.NullString{}
		} else {
			domain = sql.NullString{String: d, Valid: true}
		}
	}

	_, err = db.Exec(`UPDATE pages SET name = $1, page_id = $2, domain = $3, updated_at = NOW() WHERE id = $4 AND tenant_id = $5`, name, slug, domain, pageID, tenantID)

	if err != nil {
		if strings.Contains(err.Error(), "pages_unique_tenant_page") {
			c.JSON(409, responses.ErrorResponse{Success: false, Message: "A page with this slug already exists in this site."})
			return
		}
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Failed to update page in Postgres."})
		return
	}

	updateDoc := bson.M{"page_id": slug, "updated_at": time.Now().UTC()}

	if body.Svelte != nil {
		updateDoc["svelte"] = *body.Svelte
	}
	if body.Header != nil {
		updateDoc["header"] = *body.Header
	}
	if body.Footer != nil {
		updateDoc["footer"] = *body.Footer
	}
	if body.ShowHeader != nil {
		updateDoc["show_header"] = *body.ShowHeader
	}
	if body.ShowFooter != nil {
		updateDoc["show_footer"] = *body.ShowFooter
	}
	if body.Components != nil {
		updateDoc["components"] = body.Components
	}

	_, err = config.MongoClient.Database("genyou").Collection("page_sveltes").UpdateOne(c.Request.Context(), bson.M{"_id": pageID}, bson.M{"$set": updateDoc})

	if err != nil {
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Failed to update page content in MongoDB."})
		return
	}

	_, _ = config.MongoClient.Database("genyou").Collection("theme_store_entries").UpdateOne(c.Request.Context(), bson.M{"_id": pageID}, bson.M{"$set": bson.M{
		"name":       name,
		"page_id":    slug,
		"domain":     domain.String,
		"updated_at": time.Now().UTC(),
	}}, options.Update().SetUpsert(true))

	c.JSON(200, responses.UpdatePageDTO{
		Id:        pageID,
		Page_id:   slug,
		Name:      name,
		Tenant_id: tenantID,
	})
}
