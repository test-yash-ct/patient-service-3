package config

import (
	"os"
	"strconv"

	"github.com/healthops/patient-service/internal/obs"
)

type Config struct {
	DatabaseURL string
	ListenAddr  string
	JWTSecret   string
	Debug       bool
	Metadata    obs.Metadata
}

func Load() Config {
	debug, _ := strconv.ParseBool(getenv("DEBUG", "false"))
	return Config{
		DatabaseURL: getenv("DATABASE_URL", "postgres://app:app@localhost:5432/patients?sslmode=disable"),
		ListenAddr:  getenv("LISTEN_ADDR", "0.0.0.0:8080"),
		JWTSecret:   getenv("JWT_SECRET", "dev-only-secret"),
		Debug:       debug,
		Metadata:    obs.MetadataFromEnv("patient-service"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
