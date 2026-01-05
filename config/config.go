package config

import (
	"errors"
	"os"
	"time"
)

type Config struct {
	Port      string
	JWTSecret string
	JWTTTL    time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Port:      getenvDefault("PORT", "8080"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTTTL:    parseDurationDefault("JWT_TTL", 24*time.Hour),
	}

	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}

	return cfg, nil
}

func getenvDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseDurationDefault(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return fallback

}
