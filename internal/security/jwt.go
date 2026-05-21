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

// accessKey and refreshKey are functions (not vars) so they read the config
// globals at call time — after config.Load() has populated them.
func accessKey() []byte  { return []byte(config.Jwt_access_token) }
func refreshKey() []byte { return []byte(config.Jwt_refresh_token) }

func AccessToken(claims AccessTokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":        claims.Iss,
		"sub":        claims.Sub,
		"aud":        claims.Aud,
		"exp":        claims.Exp,
		"iat":        time.Now().Unix(),
		"role":       claims.Role,
		"website_id": claims.WebsiteId,
	})
	return token.SignedString(accessKey())
}

func RefreshToken(claims RefreshTokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":        claims.Iss,
		"sub":        claims.Sub,
		"website_id": claims.WebsiteId,
		"exp":        claims.Exp,
		"iat":        time.Now().Unix(),
		"type":       "refresh",
	})
	return token.SignedString(refreshKey())
}

func ParseAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return accessKey(), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	getString := func(key string) (string, bool) {
		v, ok := claimsMap[key].(string)
		return v, ok
	}
	getFloat := func(key string) (float64, bool) {
		v, ok := claimsMap[key].(float64)
		return v, ok
	}

	iss, ok1 := getString("iss")
	sub, ok2 := getString("sub")
	aud, ok3 := getString("aud")
	role, ok4 := getString("role")
	exp, ok5 := getFloat("exp")
	iat, ok6 := getFloat("iat")
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return nil, errors.New("malformed token claims")
	}

	websiteID, _ := getString("website_id")

	return &AccessTokenClaims{
		Iss:       iss,
		Sub:       sub,
		Aud:       aud,
		Role:      role,
		WebsiteId: websiteID,
		Exp:       int64(exp),
		Iat:       int64(iat),
	}, nil
}

func ParseRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return refreshKey(), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	if claimsMap["type"] != "refresh" {
		return nil, errors.New("not a refresh token")
	}

	getString := func(key string) (string, bool) {
		v, ok := claimsMap[key].(string)
		return v, ok
	}
	getFloat := func(key string) (float64, bool) {
		v, ok := claimsMap[key].(float64)
		return v, ok
	}

	iss, ok1 := getString("iss")
	sub, ok2 := getString("sub")
	websiteID, ok3 := getString("website_id")
	typ, ok4 := getString("type")
	exp, ok5 := getFloat("exp")
	iat, ok6 := getFloat("iat")
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return nil, errors.New("malformed token claims")
	}

	return &RefreshTokenClaims{
		Iss:       iss,
		Sub:       sub,
		WebsiteId: websiteID,
		Exp:       int64(exp),
		Iat:       int64(iat),
		Type:      typ,
	}, nil
}
