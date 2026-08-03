package token

import (
	"time"

	"github.com/o1egl/paseto/v2"
)

type Claims struct {
	UserUUID    string `json:"user_uuid"`
	WebSiteUUID string `json:"website_uuid"`
	Role        string `json:"role"`
}

type pasetoClaims struct {
	Type        string    `json:"type"`
	UserUUID    string    `json:"user_uuid"`
	WebSiteUUID string    `json:"website_uuid"`
	Role        string    `json:"role"`
	IssuedAt    time.Time `json:"iat"`
	ExpiresAt   time.Time `json:"exp"`
}

const (
	accessTTL  = 15 * time.Minute
	refreshTTL = 7 * 24 * time.Hour
)

func GenerateAccess(secret []byte, userUUID string, websiteUUID string, role string) (string, error) {
	return generate(secret, "access", userUUID, websiteUUID, role, accessTTL)
}

func GenerateRefresh(secret []byte, userUUID string, websiteUUID string, role string) (string, error) {
	return generate(secret, "refresh", userUUID, websiteUUID, role, refreshTTL)
}

func generate(secret []byte, tokenType string, userUUID string, websiteUUID string, role string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := pasetoClaims{
		Type:        tokenType,
		UserUUID:    userUUID,
		WebSiteUUID: websiteUUID,
		Role:        role,
		IssuedAt:    now,
		ExpiresAt:   now.Add(ttl),
	}

	return paseto.NewV2().Encrypt(secret, claims, nil)
}

func Parse(tokenStr string, secret []byte) (*Claims, error) {
	var claims pasetoClaims

	err := paseto.NewV2().Decrypt(tokenStr, secret, &claims, nil)
	if err != nil {
		return nil, err
	}

	if time.Now().After(claims.ExpiresAt) {
		return nil, ErrExpired
	}

	return &Claims{
		UserUUID:    claims.UserUUID,
		WebSiteUUID: claims.WebSiteUUID,
		Role:        claims.Role,
	}, nil
}

var ErrExpired = errExpired{}

type errExpired struct{}

func (e errExpired) Error() string { return "token expired" }
