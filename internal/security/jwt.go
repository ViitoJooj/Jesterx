package security

import (
	"errors"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	Iss       string
	Sub       string
	Aud       string
	WebsiteId string
	Exp       int64
	Iat       int64
	Role      string
}

type RefreshTokenClaims struct {
	Iss       string
	Sub       string
	WebsiteId string
	Exp       int64
	Iat       int64
	Type      string
}

func RefreshCookieName(websiteId string) string {
	return "refresh_token_" + websiteId
}

func AccessCookieName(websiteId string) string {
	return "access_token_" + websiteId
}

func jwtAccessTokenKey() []byte { return []byte(config.Jwt_access_token) }
func jwtRefreshTokenKey() []byte { return []byte(config.Jwt_refresh_token) }

func AccessToken(claims AccessTokenClaims) (string, error) {
	now := time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  claims.Iss,
		"sub":  claims.Sub,
		"aud":  claims.Aud,
		"exp":  claims.Exp,
		"iat":  now,
		"role": claims.Role,
	})

	return token.SignedString(jwtAccessTokenKey())
}

func RefreshToken(claims RefreshTokenClaims) (string, error) {
	now := time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":        claims.Iss,
		"sub":        claims.Sub,
		"website_id": claims.WebsiteId,
		"exp":        claims.Exp,
		"iat":        now,
		"type":       "refresh",
	})

	return token.SignedString(jwtRefreshTokenKey())
}

func ParseAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Internal error.")
		}
		return jwtAccessTokenKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token.")
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid claims.")
	}

	claims := &AccessTokenClaims{
		Iss:  claimsMap["iss"].(string),
		Sub:  claimsMap["sub"].(string),
		Aud:  claimsMap["aud"].(string),
		Role: claimsMap["role"].(string),
		Exp:  int64(claimsMap["exp"].(float64)),
		Iat:  int64(claimsMap["iat"].(float64)),
	}

	return claims, nil
}

func ParseRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Internal Error.")
		}
		return jwtRefreshTokenKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token.")
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid claims.")
	}

	if claimsMap["type"] != "refresh" {
		return nil, errors.New("is not refresh token.")
	}

	claims := &RefreshTokenClaims{
		Iss:       claimsMap["iss"].(string),
		Sub:       claimsMap["sub"].(string),
		WebsiteId: claimsMap["website_id"].(string),
		Exp:       int64(claimsMap["exp"].(float64)),
		Iat:       int64(claimsMap["iat"].(float64)),
		Type:      claimsMap["type"].(string),
	}

	return claims, nil
}
