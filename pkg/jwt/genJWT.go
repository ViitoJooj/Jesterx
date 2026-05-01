package jwt

import (
	"errors"
	"strconv"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenRefreshToken(user *domain.Profile) (string, error) {
	tokenUUID, err := uuid.NewV7()
	if err != nil {
		return "", errors.New("internal error")
	}

	expDays, err := strconv.Atoi(dotenv.RefreshTokenExpDays)
	if err != nil {
		return "", errors.New("internal error")
	}

	now := time.Now()

	claims := jwt.RegisteredClaims{
		ID:        tokenUUID.String(),
		Issuer:    "JesterX",
		Subject:   user.Uid.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expDays) * 24 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(dotenv.SecretKey))
	if err != nil {
		return "", errors.New("internal error")
	}

	return signedToken, nil
}

func GenAccessToken(user *domain.Profile) (string, error) {
	expMinutes, err := strconv.Atoi(dotenv.AccessTokenExpMinutes)
	if err != nil {
		return "", errors.New("internal error")
	}

	now := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    "JesterX",
		Subject:   user.Uid.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expMinutes) * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(dotenv.SecretKey))
	if err != nil {
		return "", errors.New("internal error")
	}

	return signedToken, nil
}
