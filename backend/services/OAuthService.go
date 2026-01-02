package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"jesterx-core/config"
	"jesterx-core/helpers"
	"jesterx-core/responses"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

type GithubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type GithubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

func GoogleLoginService(c *gin.Context) {
	if config.GoogleOAuthConfig == nil {
		c.JSON(http.StatusServiceUnavailable, responses.ErrorResponse{Success: false, Message: "Google OAuth not configured"})
		return
	}

	state := uuid.New().String()
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	url := config.GoogleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallbackService(c *gin.Context) {
	if config.GoogleOAuthConfig == nil {
		c.JSON(http.StatusServiceUnavailable, responses.ErrorResponse{Success: false, Message: "Google OAuth not configured"})
		return
	}

	state := c.Query("state")
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil || state != stateCookie {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Invalid state parameter"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Code not found"})
		return
	}

	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to exchange token"})
		return
	}

	client := config.GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to read user info"})
		return
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to parse user info"})
		return
	}

	if !userInfo.VerifiedEmail {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Email not verified"})
		return
	}

	user, err := findOrCreateOAuthUser(c.Request.Context(), userInfo.Email, userInfo.GivenName, userInfo.FamilyName, userInfo.Picture, "google")
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}

	if user.Banned {
		c.JSON(http.StatusForbidden, responses.ErrorResponse{Success: false, Message: "User is banned"})
		return
	}

	user.Role = helpers.ResolvePlatformRole(user.Email, user.Role)
	if user.Role != "" {
		_, _ = config.DB.Exec(`UPDATE users SET role = $1 WHERE id = $2`, user.Role, user.Id)
	}

	jwtToken, err := helpers.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Internal error"})
		return
	}

	helpers.SetAuthCookie(c, jwtToken, 7*24*3600)
	c.Redirect(http.StatusTemporaryRedirect, config.HostProd)
}

func GithubLoginService(c *gin.Context) {
	if config.GithubOAuthConfig == nil {
		c.JSON(http.StatusServiceUnavailable, responses.ErrorResponse{Success: false, Message: "GitHub OAuth not configured"})
		return
	}

	state := uuid.New().String()
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	url := config.GithubOAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GithubCallbackService(c *gin.Context) {
	if config.GithubOAuthConfig == nil {
		c.JSON(http.StatusServiceUnavailable, responses.ErrorResponse{Success: false, Message: "GitHub OAuth not configured"})
		return
	}

	state := c.Query("state")
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil || state != stateCookie {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Invalid state parameter"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "Code not found"})
		return
	}

	token, err := config.GithubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to exchange token"})
		return
	}

	client := config.GithubOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to read user info"})
		return
	}

	var userInfo GithubUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Failed to parse user info"})
		return
	}

	email := userInfo.Email
	if email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			emailData, err := io.ReadAll(emailResp.Body)
			if err == nil {
				var emails []GithubEmail
				if json.Unmarshal(emailData, &emails) == nil {
					for _, e := range emails {
						if e.Primary && e.Verified {
							email = e.Email
							break
						}
					}
					if email == "" && len(emails) > 0 && emails[0].Verified {
						email = emails[0].Email
					}
				}
			}
		}
	}

	if email == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Success: false, Message: "No verified email found"})
		return
	}

	firstName, lastName := splitName(userInfo.Name)
	user, err := findOrCreateOAuthUser(c.Request.Context(), email, firstName, lastName, userInfo.AvatarURL, "github")
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Database error"})
		return
	}

	if user.Banned {
		c.JSON(http.StatusForbidden, responses.ErrorResponse{Success: false, Message: "User is banned"})
		return
	}

	user.Role = helpers.ResolvePlatformRole(user.Email, user.Role)
	if user.Role != "" {
		_, _ = config.DB.Exec(`UPDATE users SET role = $1 WHERE id = $2`, user.Role, user.Id)
	}

	jwtToken, err := helpers.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Success: false, Message: "Internal error"})
		return
	}

	helpers.SetAuthCookie(c, jwtToken, 7*24*3600)
	c.Redirect(http.StatusTemporaryRedirect, config.HostProd)
}

func findOrCreateOAuthUser(ctx context.Context, email, firstName, lastName, profileImg, provider string) (helpers.UserData, error) {
	var user helpers.UserData

	err := config.DB.QueryRowContext(ctx, `
		SELECT 
			u.id,
			u.email,
			COALESCE(u.plan, 'free'),
			COALESCE(p.profile_img, ''),
			COALESCE(p.first_name, ''),
			COALESCE(p.last_name, ''),
			COALESCE(u.role, 'platform_user'),
			COALESCE(u.banned, FALSE)
		FROM users u
		LEFT JOIN user_profiles p ON p.user_id = u.id
		WHERE u.email = $1
	`, email).Scan(
		&user.Id,
		&user.Email,
		&user.Plan,
		&user.Profile_img,
		&user.First_name,
		&user.Last_name,
		&user.Role,
		&user.Banned,
	)

	if err == nil {
		return user, nil
	}

	if err != sql.ErrNoRows {
		return helpers.UserData{}, err
	}

	userID := uuid.New().String()
	randomPassword := uuid.New().String()
	hashedPassword, err := helpers.HashPassword(randomPassword)
	if err != nil {
		return helpers.UserData{}, err
	}

	_, err = config.DB.ExecContext(ctx, `
		INSERT INTO users (id, email, password, plan, role)
		VALUES ($1, $2, $3, 'free', 'platform_user')
	`, userID, email, hashedPassword)
	if err != nil {
		return helpers.UserData{}, err
	}

	_, err = config.DB.ExecContext(ctx, `
		INSERT INTO user_profiles (user_id, first_name, last_name, profile_img)
		VALUES ($1, $2, $3, $4)
	`, userID, firstName, lastName, profileImg)
	if err != nil {
		return helpers.UserData{}, err
	}

	user.Id = userID
	user.Email = email
	user.Plan = "free"
	user.Profile_img = profileImg
	user.First_name = firstName
	user.Last_name = lastName
	user.Role = "platform_user"
	user.Banned = false

	return user, nil
}

func splitName(fullName string) (string, string) {
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return "", ""
	}

	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}

	firstName := parts[0]
	lastName := strings.Join(parts[1:], " ")
	return firstName, lastName
}
