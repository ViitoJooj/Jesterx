package jwt

import (
	"errors"
	"fmt"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ValidateAccessToken(tokenValue string, cfg dotenv.JWTConfig) (*jwt.RegisteredClaims, error) {
	claims, err := validateToken(tokenValue, cfg)
	if err != nil {
		return nil, err
	}

	if claims.ID != "" {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func ValidateRefreshToken(tokenValue string, cfg dotenv.JWTConfig) (*jwt.RegisteredClaims, error) {
	claims, err := validateToken(tokenValue, cfg)
	if err != nil {
		return nil, err
	}

	if claims.ID == "" {
		return nil, errors.New("invalid token")
	}

	if _, err := uuid.Parse(claims.ID); err != nil {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func validateToken(tokenValue string, cfg dotenv.JWTConfig) (*jwt.RegisteredClaims, error) {
	if tokenValue == "" {
		return nil, errors.New("invalid token")
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Issuer != "JesterX" {
		return nil, errors.New("invalid token")
	}

	if claims.Subject == "" {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
