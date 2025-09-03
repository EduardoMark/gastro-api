package config

import (
	"os"
)

type Env struct {
	DbUser     string
	DbPassword string
	DbName     string
	DbPort     string
	DbHost     string
	JWTSecret  string
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Load() *Env {
	cfg := &Env{
		DbUser:     getEnv("DB_USER", "postgres"),
		DbPassword: getEnv("DB_PASSWORD", "root"),
		DbName:     getEnv("DB_NAME", "gastro_api"),
		DbPort:     getEnv("DB_PORT", "5432"),
		DbHost:     getEnv("DB_HOST", "db"),
		JWTSecret:  getEnv("JWT_SECRET", "secret-key"),
	}

	return cfg
}
