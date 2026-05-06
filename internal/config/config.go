// internal/config/config.go
package config

import (
	"os"
)

type Config struct {
	ServerAddr  string
	DatabaseDSN string
	JWTSecret   string
	JWTDuration string
}

func Load() *Config {
	return &Config{
		ServerAddr:  getEnv("SERVER_ADDR", ":8080"),
		DatabaseDSN: getEnv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/subscriptions_db?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "secret"),
		JWTDuration: getEnv("JWT_DURATION", "24h"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
