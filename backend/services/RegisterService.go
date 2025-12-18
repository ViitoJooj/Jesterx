package services

import (
	"database/sql"
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/helpers"
	"gen-you-ecommerce/models"
	"gen-you-ecommerce/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterService(c *gin.Context) {

	tenantValue, hasTenant := c.Get("tenantID")
	var tenantID string
	if hasTenant {
		tenantID = tenantValue.(string)
	}

	var body models.RegisterModel
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Invalid request body"})
		return
	}

	if err := helpers.ValidateEmail(body.Email); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	if err := helpers.ValidatePassword(body.Password); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	exists, err := helpers.EmailExists(c.Request.Context(), config.DB, body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Email already registered"})
		return
	}

	hashedPassword, err := helpers.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Internal error"})
		return
	}

	tx, err := config.DB.BeginTx(c.Request.Context(), &sql.TxOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}
	defer tx.Rollback()

	var userID, userPlan string
	err = tx.QueryRow(`
        INSERT INTO users (email, password)
        VALUES ($1, $2)
        RETURNING id, plan
    `, body.Email, hashedPassword).Scan(&userID, &userPlan)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	_, err = tx.Exec(`
        INSERT INTO user_profiles (user_id, first_name, last_name)
        VALUES ($1, $2, $3)
    `, userID, body.First_name, body.Last_name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	role := "platform_user"
	if hasTenant && tenantID != "" {
		_, err = tx.Exec(`
            INSERT INTO tenant_users (tenant_id, user_id, role)
            VALUES ($1, $2, 'customer')
        `, tenantID, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Success: false,
				Message: "Database error",
			})
			return
		}

		role = "customer"
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	user := helpers.UserData{
		Id:          userID,
		Email:       body.Email,
		Profile_img: "",
		First_name:  body.First_name,
		Last_name:   body.Last_name,
		Role:        role,
		Plan:        userPlan,
	}

	token, err := helpers.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Internal error"})
		return
	}

	loggedTimer := helpers.GetLoginDuration(body.Keep_me_logged_in)
	helpers.SetAuthCookie(c, token, loggedTimer)

	c.JSON(http.StatusOK, responses.UserRegisterResponse{
		Success: true,
		Message: "Registration successful.",
		Data: responses.UserData{
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
