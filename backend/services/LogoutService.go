package services

import (
	"gen-you-ecommerce/helpers"
	"gen-you-ecommerce/responses"

	"github.com/gin-gonic/gin"
)

func LogoutService(c *gin.Context) {
	helpers.SetAuthCookie(c, "", 0)

	c.JSON(200, responses.LogoutResponse{
		Success: true,
		Message: "The user successfully exited the session.",
	})
}
