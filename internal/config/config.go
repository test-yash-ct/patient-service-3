package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL    string
	ListenAddr     string
	JWTSecret      string
	PIIKey         []byte
	Debug          bool
	CORSOrigins    []string
	MaxTokenTTLSec int64
}

func Load() (Config, error) {
	debug, _ := strconv.ParseBool(getenv("DEBUG", "false"))
	secret := os.Getenv("JWT_SECRET")
	if secret == "" || secret == "dev-only-secret" {
		return Config{}, fmt.Errorf("JWT_SECRET must be set to a non-default value")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if strings.Contains(strings.ToLower(dbURL), "sslmode=disable") {
		return Config{}, fmt.Errorf("DATABASE_URL must use TLS (sslmode=require or verify-full)")
	}
	key, err := loadPIIKey()
	if err != nil {
		return Config{}, err
	}
	origins := splitCSV(getenv("CORS_ORIGINS", ""))
	return Config{
		DatabaseURL:    dbURL,
		ListenAddr:     getenv("LISTEN_ADDR", "127.0.0.1:8080"),
		JWTSecret:      secret,
		PIIKey:         key,
		Debug:          debug,
		CORSOrigins:    origins,
		MaxTokenTTLSec: 900,
	}, nil
}

func loadPIIKey() ([]byte, error) {
	raw := os.Getenv("PII_ENCRYPTION_KEY")
	if len(raw) != 32 {
		return nil, fmt.Errorf("PII_ENCRYPTION_KEY must be exactly 32 bytes")
	}
	return []byte(raw), nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
