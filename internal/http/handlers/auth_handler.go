package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ViitoJooj/Jesterx/internal/http/dtos"
	"github.com/ViitoJooj/Jesterx/internal/repository"
	"github.com/ViitoJooj/Jesterx/internal/service"
	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(req.Name, req.Email, req.Password, req.Role, req.Cpf)
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyInUse) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         user.Uuid,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.Created_at,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	accessExp, _ := strconv.Atoi(dotenv.AccessTokenExpMinutes)
	refreshExp, _ := strconv.Atoi(dotenv.RefreshTokenExpDays)

	secure := dotenv.Environment == "production"

	c.SetCookie("access_token", accessToken, accessExp*60, "/", "", secure, true)
	c.SetCookie("refresh_token", refreshToken, refreshExp*24*60*60, "/api/v1/auth", "", secure, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged in"})
}

func (h *AuthHandler) Token(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	accessToken, err := h.authService.Token(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired, please login again"})
		return
	}

	accessExp, _ := strconv.Atoi(dotenv.AccessTokenExpMinutes)
	secure := dotenv.Environment == "production"

	c.SetCookie("access_token", accessToken, accessExp*60, "/", "", secure, true)
	c.JSON(http.StatusOK, gin.H{"message": "token refreshed"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")

	_ = h.authService.Logout(refreshToken)

	secure := dotenv.Environment == "production"

	c.SetCookie("access_token", "", -1, "/", "", secure, true)
	c.SetCookie("refresh_token", "", -1, "/api/v1/auth", "", secure, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
