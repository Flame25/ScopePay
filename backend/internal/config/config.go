package config

import "os"

type Config struct {
	Port          string
	DevMode       bool
	DBDriver      string
	DBDSN         string
	PublicBaseURL string
}

func Load() Config {
	return Config{
		Port:          env("PORT", "8080"),
		DevMode:       env("DEV_MODE", "1") == "1",
		DBDriver:      env("DB_DRIVER", "sqlite"),
		DBDSN:         env("DB_DSN", "scope_pay.db"),
		PublicBaseURL: env("PUBLIC_BASE_URL", "http://localhost:8080"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
