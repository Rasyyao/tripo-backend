package config

import "os"

type Config struct {
	AppPort string
}

func Load() *Config {
	return &Config{
		AppPort: getEnv("APP_PORT", "3000"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
