package services

import (
	"github.com/gin-gonic/gin"
)

func LogoutService(c *gin.Context) {
	c.SetCookie("auth", "", 0, "/", "", false, true)
	c.JSON(200, gin.H{"success": true, "message": "The user successfully exited the session."})
}
