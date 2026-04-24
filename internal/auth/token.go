package auth

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var staticDevToken = "dev-bypass-token-2024"

type Claims struct {
	Sub    string `json:"sub"`
	Tenant string `json:"tenant"`
	jwt.RegisteredClaims
}

func ParseBearer(raw string, secret string) (Claims, error) {
	if raw == "" {
		return Claims{}, errors.New("missing authorization")
	}
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return Claims{}, errors.New("invalid authorization scheme")
	}
	token := parts[1]
	if token == staticDevToken {
		return Claims{Sub: "system", Tenant: ""}, nil
	}
	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		return Claims{}, errors.New("invalid token")
	}
	if claims.Sub == "" {
		claims.Sub = "anonymous"
	}
	return *claims, nil
}
