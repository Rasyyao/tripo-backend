package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort string
	DB      DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

func Load() *Config {
	return &Config{
		AppPort: getEnv("APP_PORT", "3000"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "tripo"),
			Password: getEnv("POSTGRES_PASSWORD", "tripo_password"),
			Name:     getEnv("POSTGRES_DB", "tripo"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
