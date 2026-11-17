package config

import (
	"errors"
	"os"
)

type Config struct {
	Port          string
	DevMode       bool
	DBDriver      string
	DBDSN         string
	PublicBaseURL string
	JWTSecret     string
}

func Validate(cfg Config) error {

	if !cfg.DevMode && cfg.JWTSecret == "" {
		return errors.New("JWT_SECRET required when DEV_MODE=0")

	}

	return nil
}

func Load() Config {
	cfg := Config{
		Port:          env("PORT", "8080"),
		DevMode:       env("DEV_MODE", "1") == "1",
		DBDriver:      env("DB_DRIVER", "sqlite"),
		DBDSN:         env("DB_DSN", "scope_pay.db"),
		PublicBaseURL: env("PUBLIC_BASE_URL", "http://localhost:8080"),
		JWTSecret:     env("JWT_SECRET", ""),
	}

	if err := Validate(cfg); err != nil {
		panic(err)
	}

	return cfg
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
