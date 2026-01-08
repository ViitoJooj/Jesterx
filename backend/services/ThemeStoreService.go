package services

import (
	"database/sql"
	"jesterx-core/config"
	"jesterx-core/helpers"
	"jesterx-core/models"
	"jesterx-core/responses"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListThemeStoreService(c *gin.Context) {
	tenantID, _ := c.Get("tenantID")
	currentTenant := ""
	if tenantID != nil {
		currentTenant = tenantID.(string)
	}

	cursor, err := config.MongoClient.Database("genyou").Collection("theme_store_entries").Find(c.Request.Context(), bson.M{})
	if err != nil {
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Não foi possível carregar a loja de temas."})
		return
	}
	defer cursor.Close(c.Request.Context())

	type ThemeEntryResponse struct {
		ID        string    `json:"id"`
		PageID    string    `json:"page_id"`
		Name      string    `json:"name"`
		Domain    string    `json:"domain"`
		ForSale   bool      `json:"for_sale"`
		Owned     bool      `json:"owned"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	var resp []ThemeEntryResponse
	for cursor.Next(c.Request.Context()) {
		var entry models.ThemeStoreEntry
		if err := cursor.Decode(&entry); err != nil {
			continue
		}
		resp = append(resp, ThemeEntryResponse{
			ID:        entry.ID,
			PageID:    entry.PageID,
			Name:      entry.Name,
			Domain:    entry.Domain,
			ForSale:   entry.ForSale,
			Owned:     currentTenant != "" && entry.TenantID == currentTenant,
			UpdatedAt: entry.UpdatedAt,
		})
	}

	c.JSON(200, gin.H{"success": true, "data": resp})
}

func GetThemeStoreBySlugService(c *gin.Context) {
	slug := c.Param("slug")

	var entry models.ThemeStoreEntry
	err := config.MongoClient.Database("genyou").Collection("theme_store_entries").FindOne(
		c.Request.Context(),
		bson.M{"page_id": slug, "for_sale": true},
	).Decode(&entry)

	if err != nil {
		c.JSON(404, responses.ErrorResponse{Success: false, Message: "Tema não encontrado."})
		return
	}

	type ThemeDetailResponse struct {
		ID              string   `json:"id"`
		Name            string   `json:"name"`
		Description     string   `json:"description"`
		Images          []string `json:"images"`
		Rating          float64  `json:"rating"`
		Installs        int      `json:"installs"`
		LongDescription string   `json:"long_description"`
		PageID          string   `json:"page_id"`
		Domain          string   `json:"domain"`
	}

	// For now, we'll return mock data for images, rating, and installs
	// In a real scenario, this would come from the database
	resp := ThemeDetailResponse{
		ID:              entry.ID,
		Name:            entry.Name,
		Description:     "Tema moderno e responsivo para sua loja online",
		Images:          []string{"https://via.placeholder.com/800x600?text=" + entry.Name},
		Rating:          4.5,
		Installs:        150,
		LongDescription: "Este tema oferece um design limpo e profissional, perfeito para qualquer tipo de negócio. Com recursos avançados de personalização e otimizado para conversão.",
		PageID:          entry.PageID,
		Domain:          entry.Domain,
	}

	c.JSON(200, gin.H{"success": true, "data": resp})
}

func UpdateThemeStoreEntryService(c *gin.Context) {
	user := c.MustGet("user").(helpers.UserData)
	tenantID := c.MustGet("tenantID").(string)
	pageSlug := c.Param("page_id")

	var body struct {
		ForSale bool `json:"for_sale"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Dados inválidos."})
		return
	}

	hasAccess, _, err := helpers.UserHasTenantAccess(user.Id, tenantID)
	if err != nil {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Falha ao validar acesso."})
		return
	}
	if !hasAccess {
		c.JSON(403, responses.ErrorResponse{Success: false, Message: "Você não tem permissão para editar esta loja."})
		return
	}

	var pageID, name, domain string
	if err := config.DB.QueryRow(`SELECT id, name, COALESCE(domain, '') FROM pages WHERE tenant_id = $1 AND page_id = $2`, tenantID, pageSlug).Scan(&pageID, &name, &domain); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, responses.ErrorResponse{Success: false, Message: "Página não encontrada."})
			return
		}
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Erro ao validar página."})
		return
	}

	_, err = config.MongoClient.Database("genyou").Collection("theme_store_entries").UpdateOne(
		c.Request.Context(),
		bson.M{"_id": pageID},
		bson.M{"$set": bson.M{
			"for_sale":   body.ForSale,
			"name":       name,
			"page_id":    pageSlug,
			"domain":     domain,
			"tenant_id":  tenantID,
			"updated_at": time.Now().UTC(),
		}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Não foi possível atualizar a loja de temas."})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Tema atualizado na loja."})
}
