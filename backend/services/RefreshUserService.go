package services

import (
	"database/sql"
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/helpers"
	"gen-you-ecommerce/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RefreshUserService(c *gin.Context) {
	userAny, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse{Success: false, Message: "unauthorized"})
		return
	}
	user, ok := userAny.(helpers.UserData)
	if !ok || user.Id == "" {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse{Success: false, Message: "unauthorized"})
		return
	}

	var updatedUser helpers.UserData
	err := config.DB.QueryRow(`
		SELECT 
			u.id,
			u.email,
			COALESCE(u.plan, 'free'),
			COALESCE(p.profile_img, ''),
			COALESCE(p.first_name, ''),
			COALESCE(p. last_name, '')
		FROM users u
		LEFT JOIN user_profiles p ON p.user_id = u.id
		WHERE u.id = $1
	`, user.Id).Scan(
		&updatedUser.Id,
		&updatedUser.Email,
		&updatedUser.Plan,
		&updatedUser.Profile_img,
		&updatedUser.First_name,
		&updatedUser.Last_name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, responses.ErrorResponse{Success: false, Message: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}

	updatedUser.Role = user.Role

	token, err := helpers.GenerateToken(updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Token generation failed"})
		return
	}

	helpers.SetAuthCookie(c, token, 24*7)

	c.JSON(http.StatusOK, responses.UserLoginResponse{
		Success: true,
		Message: "User data refreshed",
		Data: responses.UserData{
			Id:         updatedUser.Id,
			ProfileImg: updatedUser.Profile_img,
			FirstName:  updatedUser.First_name,
			LastName:   updatedUser.Last_name,
			Email:      updatedUser.Email,
			Role:       updatedUser.Role,
			Plan:       updatedUser.Plan,
		},
	})
}
