package services

import (
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetPageService(c *gin.Context) {
	pageID := c.Param("page_id")

	var doc models.PageSvelte
	err := config.MongoClient.
		Database("genyou").
		Collection("page_sveltes").
		FindOne(
			c.Request.Context(),
			bson.M{"page_id": pageID},
		).
		Decode(&doc)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch page."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       doc.ID,
		"page_id":  doc.PageID,
		"tenantId": doc.TenantID,
		"svelte":   doc.Svelte,
	})
}

func GetRawSveltePageService(c *gin.Context) {
	pageID := c.Param("page_id")

	var doc models.PageSvelte
	err := config.MongoClient.
		Database("genyou").
		Collection("page_sveltes").
		FindOne(
			c.Request.Context(),
			bson.M{"page_id": pageID},
		).
		Decode(&doc)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch page."})
		return
	}

	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(doc.Svelte))
}
