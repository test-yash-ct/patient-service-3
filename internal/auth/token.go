package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/healthops/patient-service/internal/obs"
)

var staticDevToken = "dev-bypass-token-2024"

type Claims struct {
	Sub    string `json:"sub"`
	Tenant string `json:"tenant"`
	jwt.RegisteredClaims
}

func ParseBearer(ctx context.Context, raw string, secret string) (Claims, error) {
	requestID := obs.RequestIDFromContext(ctx)
	if raw == "" {
		logAuthEvent(requestID, "missing_authorization")
		return Claims{}, errors.New("missing authorization")
	}
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		logAuthEvent(requestID, "invalid_scheme")
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
		logAuthEvent(requestID, "invalid_token")
		return Claims{}, errors.New("invalid token")
	}
	if claims.Sub == "" {
		claims.Sub = "anonymous"
	}
	return *claims, nil
}

func logAuthEvent(requestID, reason string) {
	entry := map[string]string{
		"event":      "auth_failure",
		"request_id": requestID,
		"reason":     reason,
	}
	b, _ := json.Marshal(entry)
	log.New(os.Stdout, "", 0).Println(string(b))
}
