package configs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"task_tracker/internal/helpers"
)

type AuthConfig struct {
	JWTSecret            string
	AccessTokenTTL       time.Duration
	LegacyHeadersEnabled bool
}

func LoadAuth() (AuthConfig, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if isProduction() && jwtSecret == "" {
		return AuthConfig{}, fmt.Errorf("JWT_SECRET must be set in production")
	}
	if jwtSecret == "" {
		jwtSecret = "dev-only-change-me-at-least-32-bytes"
	}

	return AuthConfig{
		JWTSecret:            jwtSecret,
		AccessTokenTTL:       loadDurationMinutes("JWT_ACCESS_TTL_MINUTES", 15),
		LegacyHeadersEnabled: loadBool("AUTH_LEGACY_HEADERS_ENABLED", defaultLegacyHeadersEnabled()),
	}, nil
}

func loadDurationMinutes(key string, defaultValue int) time.Duration {
	raw := helpers.GetEnv(key, "")
	if raw == "" {
		return time.Duration(defaultValue) * time.Minute
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return time.Duration(defaultValue) * time.Minute
	}
	return time.Duration(value) * time.Minute
}

func loadBool(key string, defaultValue bool) bool {
	raw := helpers.GetEnv(key, "")
	if raw == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return defaultValue
	}
	return value
}

func defaultLegacyHeadersEnabled() bool {
	return !isProduction()
}

func isProduction() bool {
	appEnv := strings.ToLower(helpers.GetEnv("APP_ENV", helpers.GetEnv("ENV", "development")))
	ginMode := strings.ToLower(helpers.GetEnv("GIN_MODE", ""))
	return appEnv == "production" || ginMode == "release"
}
