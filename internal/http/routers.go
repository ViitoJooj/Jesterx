package http

import "github.com/gin-gonic/gin"

func RegisterRouters() *gin.Engine {
	router := gin.Default()
	return router
}
