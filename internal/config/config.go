package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const minJWTSecretLength = 32

type Config struct {
	AppPort string
	DB      DBConfig
	JWT     JWTConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig controls signing and lifetimes of access and refresh tokens.
type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
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
		JWT: JWTConfig{
			Secret:     os.Getenv("JWT_SECRET"),
			AccessTTL:  getDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTTL: getDuration("JWT_REFRESH_TTL", 30*24*time.Hour),
		},
	}
}

// Validate reports configuration that would make the server unsafe or unable to start.
func (c *Config) Validate() error {
	if len(c.JWT.Secret) < minJWTSecretLength {
		return fmt.Errorf("JWT_SECRET must be set and at least %d characters", minJWTSecretLength)
	}
	if c.JWT.AccessTTL <= 0 || c.JWT.RefreshTTL <= 0 {
		return errors.New("JWT_ACCESS_TTL and JWT_REFRESH_TTL must be positive durations")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return d
}
