package services

import (
	"database/sql"
	"jesterx-core/config"
	"jesterx-core/helpers"
	"jesterx-core/models"
	"jesterx-core/responses"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateProductService(c *gin.Context) {
	user := c.MustGet("user").(helpers.UserData)
	tenantID := c.MustGet("tenantID").(string)
	pageSlug := c.Param("page_id")

	if strings.TrimSpace(pageSlug) == "" {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Página inválida para produto."})
		return
	}

	var body struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		PriceCents  int64    `json:"price_cents"`
		Images      []string `json:"images"`
		Visible     *bool    `json:"visible"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Dados de produto inválidos."})
		return
	}

	hasAccess, role, err := helpers.UserHasTenantAccess(user.Id, tenantID)
	if err != nil {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Falha ao validar acesso."})
		return
	}
	if !hasAccess || (role != "owner" && role != "admin" && role != "editor") {
		c.JSON(403, responses.ErrorResponse{Success: false, Message: "Você não tem permissão para gerenciar produtos."})
		return
	}

	var pageID string
	if err := config.DB.QueryRow(`SELECT id FROM pages WHERE tenant_id = $1 AND page_id = $2`, tenantID, pageSlug).Scan(&pageID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, responses.ErrorResponse{Success: false, Message: "Página não encontrada para este tenant."})
			return
		}
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Erro ao validar página."})
		return
	}

	visible := true
	if body.Visible != nil {
		visible = *body.Visible
	}

	now := time.Now().UTC()
	product := models.Product{
		ID:          uuid.NewString(),
		TenantID:    tenantID,
		PageID:      pageSlug,
		Name:        strings.TrimSpace(body.Name),
		Description: strings.TrimSpace(body.Description),
		PriceCents:  body.PriceCents,
		Images:      body.Images,
		Visible:     visible,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err = config.MongoClient.Database("genyou").Collection("products").InsertOne(c.Request.Context(), product)
	if err != nil {
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Não foi possível salvar o produto."})
		return
	}

	c.JSON(201, gin.H{"success": true, "data": product})
}

func UpdateProductService(c *gin.Context) {
	user := c.MustGet("user").(helpers.UserData)
	tenantID := c.MustGet("tenantID").(string)
	pageSlug := c.Param("page_id")
	productID := c.Param("product_id")

	var body struct {
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		PriceCents  *int64   `json:"price_cents"`
		Images      []string `json:"images"`
		Visible     *bool    `json:"visible"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Dados de produto inválidos."})
		return
	}

	hasAccess, role, err := helpers.UserHasTenantAccess(user.Id, tenantID)
	if err != nil {
		c.JSON(400, responses.ErrorResponse{Success: false, Message: "Falha ao validar acesso."})
		return
	}
	if !hasAccess || (role != "owner" && role != "admin" && role != "editor") {
		c.JSON(403, responses.ErrorResponse{Success: false, Message: "Você não tem permissão para editar produtos."})
		return
	}

	// Ensure page exists for tenant
	if err := config.DB.QueryRow(`SELECT id FROM pages WHERE tenant_id = $1 AND page_id = $2`, tenantID, pageSlug).Scan(new(string)); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, responses.ErrorResponse{Success: false, Message: "Página não encontrada para este tenant."})
			return
		}
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Erro ao validar página."})
		return
	}

	update := bson.M{"updated_at": time.Now().UTC()}
	if body.Name != nil {
		update["name"] = strings.TrimSpace(*body.Name)
	}
	if body.Description != nil {
		update["description"] = strings.TrimSpace(*body.Description)
	}
	if body.PriceCents != nil {
		update["price_cents"] = *body.PriceCents
	}
	if body.Images != nil {
		update["images"] = body.Images
	}
	if body.Visible != nil {
		update["visible"] = *body.Visible
	}

	result, err := config.MongoClient.Database("genyou").Collection("products").UpdateOne(
		c.Request.Context(),
		bson.M{"_id": productID, "tenant_id": tenantID, "page_id": pageSlug},
		bson.M{"$set": update},
	)
	if err != nil {
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Erro ao atualizar produto."})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(404, responses.ErrorResponse{Success: false, Message: "Produto não encontrado."})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Produto atualizado."})
}

func ListProductsService(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(string)
	pageSlug := c.Param("page_id")

	products := []models.Product{}
	cursor, err := config.MongoClient.Database("genyou").Collection("products").Find(c.Request.Context(), bson.M{
		"tenant_id": tenantID,
		"page_id":   pageSlug,
	})
	if err != nil {
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Não foi possível carregar produtos."})
		return
	}
	defer cursor.Close(c.Request.Context())

	for cursor.Next(c.Request.Context()) {
		var p models.Product
		if err := cursor.Decode(&p); err == nil {
			products = append(products, p)
		}
	}

	c.JSON(200, gin.H{"success": true, "data": products})
}

func PublicListProductsService(c *gin.Context) {
	tenantID := c.MustGet("tenantID").(string)
	pageSlug := c.Param("page_id")

	products := []models.Product{}
	cursor, err := config.MongoClient.Database("genyou").Collection("products").Find(c.Request.Context(), bson.M{
		"tenant_id": tenantID,
		"page_id":   pageSlug,
		"visible":   true,
	})
	if err != nil {
		c.JSON(500, responses.ErrorResponse{Success: false, Message: "Não foi possível carregar produtos públicos."})
		return
	}
	defer cursor.Close(c.Request.Context())

	for cursor.Next(c.Request.Context()) {
		var p models.Product
		if err := cursor.Decode(&p); err == nil {
			products = append(products, p)
		}
	}

	c.JSON(200, gin.H{"success": true, "data": products})
}
