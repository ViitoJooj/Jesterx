package jwt

import (
	"errors"
	"fmt"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	"github.com/ViitoJooj/Jesterx/pkg/validators"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ValidateAccessToken(tokenValue string) (*jwt.RegisteredClaims, error) {
	claims, err := validateToken(tokenValue)
	if err != nil {
		return nil, err
	}

	if claims.ID != "" {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func ValidateRefreshToken(tokenValue string) (*jwt.RegisteredClaims, error) {
	claims, err := validateToken(tokenValue)
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

func validateToken(tokenValue string) (*jwt.RegisteredClaims, error) {
	if tokenValue == "" {
		return nil, errors.New("invalid token")
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(dotenv.SecretKey), nil
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

	if err := validators.Uuid(claims.Subject); err != nil {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
