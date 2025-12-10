package services

import (
	"gen-you-ecommerce/helpers"
	"gen-you-ecommerce/responses"

	"github.com/gin-gonic/gin"
)

func MeService(c *gin.Context) {
	user := c.MustGet("user").(helpers.UserData)

	c.JSON(200, responses.MeResponse{
		Success: true,
		Message: "User info",
		Data: responses.UserDataMe{
			Id:         user.Id,
			ProfileImg: user.Profile_img,
			FirstName:  user.First_name,
			LastName:   user.Last_name,
			Email:      user.Email,
			Role:       user.Role,
			Plan:       user.Plan,
		},
	})
}
