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

func LoginService(c *gin.Context) {
	var body models.LoginModel
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
	if !exists {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Verify email or password"})
		return
	}

	tenantIDValue, hasTenant := c.Get("tenantID")
	if !hasTenant || tenantIDValue == nil {
		loginPlatform(c, body)
		return
	}

	tenantID, ok := tenantIDValue.(string)
	if !ok || tenantID == "" {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Invalid tenant context"})
		return
	}

	loginTenant(c, body, tenantID)
}

func loginPlatform(c *gin.Context, body models.LoginModel) {
	var user helpers.UserData
	var hashedPassword string

	err := config.DB.QueryRow(`
        SELECT 
            u.id,
            u.email,
            u.password,
            COALESCE(u.plan, 'free'),
            COALESCE(p.profile_img, ''),
            COALESCE(p.first_name, ''),
            COALESCE(p.last_name, '')
        FROM users u
        LEFT JOIN user_profiles p ON p.user_id = u.id
        WHERE u.email = $1
    `, body.Email).Scan(
		&user.Id,
		&user.Email,
		&hashedPassword,
		&user.Plan,
		&user.Profile_img,
		&user.First_name,
		&user.Last_name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Verify email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}

	if !helpers.CheckPasswordHash(body.Password, hashedPassword) {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Verify email or password"})
		return
	}

	user.Role = "platform_user"

	token, err := helpers.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Internal error"})
		return
	}

	loggedTimer := helpers.GetLoginDuration(body.Keep_me_logged_in)
	helpers.SetAuthCookie(c, token, loggedTimer)

	c.JSON(http.StatusOK, responses.UserLoginResponse{
		Success: true,
		Message: "Login successful.",
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

func loginTenant(c *gin.Context, body models.LoginModel, tenantID string) {
	var user helpers.UserData
	var hashedPassword string

	err := config.DB.QueryRow(`
        SELECT 
            u.id,
            u.email,
            u.password,
            COALESCE(u.plan, 'free'),
            COALESCE(p.profile_img, ''),
            COALESCE(p.first_name, ''),
            COALESCE(p.last_name, ''),
            COALESCE(tu.role, 'customer')
        FROM users u
        LEFT JOIN user_profiles p ON p.user_id = u.id
        LEFT JOIN tenant_users tu ON tu.user_id = u.id AND tu.tenant_id = $2
        WHERE u.email = $1
    `, body.Email, tenantID).Scan(
		&user.Id,
		&user.Email,
		&hashedPassword,
		&user.Plan,
		&user.Profile_img,
		&user.First_name,
		&user.Last_name,
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Verify email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}

	var hasAccess bool
	err = config.DB.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM tenant_users WHERE user_id = $1 AND tenant_id = $2)`,
		user.Id, tenantID,
	).Scan(&hasAccess)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}

	if !hasAccess || user.Role == "" {
		c.JSON(http.StatusForbidden, responses.ErrorResponse{Success: false, Message: "User does not belong to this tenant."})
		return
	}

	if !helpers.CheckPasswordHash(body.Password, hashedPassword) {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Verify email or password"})
		return
	}

	token, err := helpers.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Internal error"})
		return
	}

	loggedTimer := helpers.GetLoginDuration(body.Keep_me_logged_in)
	helpers.SetAuthCookie(c, token, loggedTimer)

	c.JSON(http.StatusOK, responses.UserLoginResponse{
		Success: true,
		Message: "Login successful.",
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
