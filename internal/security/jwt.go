package security

import (
	"time"

	"github.com/ViitoJooj/Jesterx/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	Iss   string
	Sub   string
	Aud   string
	Exp   int64
	Iat   int64
	Roles []string
}

type RefreshTokenClaims struct {
	Iss  string
	Sub  string
	Exp  int64
	Iat  int64
	Jti  string
	Type string
}

var jwtAccessTokenKey = []byte(config.Jwt_access_token)
var jwtRefreshTokenKey = []byte(config.Jwt_refresh_token)

func AccessToken(claims AccessTokenClaims) (string, error) {
	now := time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   claims.Iss,
		"sub":   claims.Sub,
		"aud":   claims.Aud,
		"exp":   claims.Exp,
		"iat":   now,
		"roles": claims.Roles,
	})

	return token.SignedString(jwtAccessTokenKey)
}

func RefreshToken(claims RefreshTokenClaims) (string, error) {
	now := time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  claims.Iss,
		"sub":  claims.Sub,
		"exp":  claims.Exp,
		"iat":  now,
		"jti":  claims.Jti,
		"type": "refresh",
	})

	return token.SignedString(jwtRefreshTokenKey)
}
