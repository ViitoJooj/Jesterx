package helpers

import (
	"errors"
	"jesterx-core/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserData struct {
	Id          string
	Profile_img string
	First_name  string
	Last_name   string
	Email       string
	Role        string
	Plan        string
	Banned      bool
}

func GenerateToken(user UserData) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub":         user.Id,
		"profile_img": user.Profile_img,
		"first_name":  user.First_name,
		"last_name":   user.Last_name,
		"email":       user.Email,
		"role":        user.Role,
		"exp":         now.Add(24 * time.Hour).Unix(),
		"iat":         now.Unix(),
		"plan":        user.Plan,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JwtSecret))
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signature method")
		}
		return []byte(config.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func SetAuthCookie(c *gin.Context, token string, time int) {
	c.SetCookie(
		"auth_token",
		token,
		60*60*time,
		"/",
		"",
		true,
		true,
	)
}

func GetLoginDuration(keepMeLoggedIn bool) int {
	if !keepMeLoggedIn {
		return 24
	}
	return 744
}
