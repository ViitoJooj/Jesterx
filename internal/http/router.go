package httpx

import (
	"github.com/ViitoJooj/Jesterx/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func NewRouter(authHandler *handlers.AuthHandler) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/token", authHandler.Token)
			auth.POST("/logout", authHandler.Logout)
		}
	}

	return router
}
