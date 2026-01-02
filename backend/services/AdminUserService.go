package services

import (
	"context"
	"database/sql"
	"errors"
	"jesterx-core/config"
	"jesterx-core/helpers"
	"jesterx-core/responses"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type updateUserBody struct {
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	ProfileImg *string `json:"profile_img"`
	Plan       *string `json:"plan"`
	Role       *string `json:"role"`
}

type banUserBody struct {
	Banned bool `json:"banned"`
}

func AdminListUsersService(c *gin.Context) {
	limit := 200
	if v := c.Query("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 && parsed <= 1000 {
			limit = parsed
		}
	}

	rows, err := config.DB.QueryContext(c.Request.Context(), `
		SELECT 
			u.id,
			u.email,
			COALESCE(p.first_name, ''),
			COALESCE(p.last_name, ''),
			COALESCE(p.profile_img, ''),
			COALESCE(u.plan, 'free'),
			COALESCE(u.role, 'platform_user'),
			COALESCE(u.banned, FALSE),
			u.created_at,
			u.updated_at
		FROM users u
		LEFT JOIN user_profiles p ON p.user_id = u.id
		ORDER BY u.created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}
	defer rows.Close()

	var users []responses.AdminUserResponse
	for rows.Next() {
		var user responses.AdminUserResponse
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.ProfileImg,
			&user.Plan,
			&user.Role,
			&user.Banned,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

func AdminUpdateUserService(c *gin.Context) {
	userID := c.Param("user_id")
	if strings.TrimSpace(userID) == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "User id required"})
		return
	}

	var body updateUserBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Invalid body"})
		return
	}

	if body.Role != nil {
		role := strings.TrimSpace(*body.Role)
		allowedRoles := map[string]bool{
			"platform_user":  true,
			"platform_admin": true,
			"customer":       true,
			"admin":          true,
			"owner":          true,
		}
		if !allowedRoles[role] {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Invalid role"})
			return
		}
	}

	if body.Plan != nil {
		if _, err := GetPlanConfig(c.Request.Context(), strings.TrimSpace(*body.Plan)); err != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Invalid plan"})
			return
		}
	}

	var currentEmail string
	err := config.DB.QueryRowContext(c.Request.Context(), `SELECT email FROM users WHERE id = $1`, userID).Scan(&currentEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, responses.ErrorResponse{Success: false, Message: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}

	roleValue := ""
	if body.Role != nil {
		roleValue = helpers.ResolvePlatformRole(currentEmail, strings.TrimSpace(*body.Role))
		if roleValue == "" {
			roleValue = "platform_user"
		}
	}

	tx, err := config.DB.BeginTx(c.Request.Context(), &sql.TxOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}
	defer tx.Rollback()

	if body.FirstName != nil || body.LastName != nil || body.ProfileImg != nil {
		_, err = tx.ExecContext(c.Request.Context(), `
			INSERT INTO user_profiles (user_id, first_name, last_name, profile_img)
			VALUES ($1, COALESCE($2, ''), COALESCE($3, ''), COALESCE($4, ''))
			ON CONFLICT (user_id) DO UPDATE
			SET first_name = COALESCE($2, user_profiles.first_name),
				last_name = COALESCE($3, user_profiles.last_name),
				profile_img = COALESCE($4, user_profiles.profile_img),
				updated_at = NOW()
		`, userID, body.FirstName, body.LastName, body.ProfileImg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to update profile"})
			return
		}
	}

	if body.Plan != nil || body.Role != nil {
		_, err = tx.ExecContext(c.Request.Context(), `
			UPDATE users
			SET 
				plan = COALESCE($2, plan),
				role = COALESCE(NULLIF($3, ''), role),
				updated_at = NOW()
			WHERE id = $1
		`, userID, body.Plan, roleValue)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to update user"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}

	if roleValue != "" {
		_, _ = config.DB.ExecContext(c.Request.Context(), `UPDATE users SET role = $1 WHERE id = $2`, roleValue, userID)
	}

	user, err := fetchAdminUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, responses.ErrorResponse{Success: false, Message: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

func AdminBanUserService(c *gin.Context) {
	userID := c.Param("user_id")
	if strings.TrimSpace(userID) == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "User id required"})
		return
	}

	var body banUserBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Invalid body"})
		return
	}

	currentUser := c.MustGet("user").(helpers.UserData)
	if currentUser.Id == userID && body.Banned {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "You cannot ban yourself"})
		return
	}

	res, err := config.DB.ExecContext(c.Request.Context(), `
		UPDATE users
		SET banned = $1, updated_at = NOW()
		WHERE id = $2
	`, body.Banned, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		c.JSON(http.StatusNotFound, responses.ErrorResponse{Success: false, Message: "User not found"})
		return
	}

	user, err := fetchAdminUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"banned": body.Banned}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

func AdminDeleteUserService(c *gin.Context) {
	userID := c.Param("user_id")
	if strings.TrimSpace(userID) == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "User id required"})
		return
	}

	currentUser := c.MustGet("user").(helpers.UserData)
	if currentUser.Id == userID {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "You cannot delete your own account"})
		return
	}

	res, err := config.DB.ExecContext(c.Request.Context(), `DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		c.JSON(http.StatusNotFound, responses.ErrorResponse{Success: false, Message: "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted",
	})
}

func fetchAdminUser(ctx context.Context, userID string) (responses.AdminUserResponse, error) {
	var user responses.AdminUserResponse
	err := config.DB.QueryRowContext(ctx, `
		SELECT 
			u.id,
			u.email,
			COALESCE(p.first_name, ''),
			COALESCE(p.last_name, ''),
			COALESCE(p.profile_img, ''),
			COALESCE(u.plan, 'free'),
			COALESCE(u.role, 'platform_user'),
			COALESCE(u.banned, FALSE),
			u.created_at,
			u.updated_at
		FROM users u
		LEFT JOIN user_profiles p ON p.user_id = u.id
		WHERE u.id = $1
	`, userID).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.ProfileImg,
		&user.Plan,
		&user.Role,
		&user.Banned,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}
