package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
}

type ServerConfig struct {
	Port           string
	Mode           string // "debug" or "release"
	AllowedOrigins string
	AppEnv         string // "dev", "prod", "test"
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret               string
	TTLHour              int // Deprecated: use AccessTokenDuration
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackURL  string
}

var AppConfig *Config

func loadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port:           getEnv("PORT", "8080"),
			Mode:           getEnv("GIN_MODE", "debug"),
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
			AppEnv:         getEnv("APP_ENV", "dev"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "skeleton_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:               getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
			TTLHour:              getEnvInt("JWT_TTL_HOUR", 24),
			AccessTokenDuration:  getEnvDuration("JWT_ACCESS_TOKEN_DURATION", 10*time.Minute, "JWT_ACCESS_TOKEN_TTL"),
			RefreshTokenDuration: getEnvDuration("JWT_REFRESH_TOKEN_DURATION", 7*24*time.Hour, "JWT_REFRESH_TOKEN_TTL"),
		},
		OAuth: OAuthConfig{
			GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			GoogleCallbackURL:  getEnv("GOOGLE_CALLBACK_URL", ""),
		},
	}
}

func (c *DatabaseConfig) ConnectionString() string {
	sslMode := c.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, sslMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func getEnvDuration(key string, defaultValue time.Duration, fallbackKeys ...string) time.Duration {
	keys := append([]string{key}, fallbackKeys...)
	for _, k := range keys {
		val := strings.TrimSpace(os.Getenv(k))
		if val == "" {
			continue
		}

		valLower := strings.ToLower(val)
		if strings.HasSuffix(valLower, "days") {
			if days, err := strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(valLower, "days"))); err == nil {
				return time.Duration(days) * 24 * time.Hour
			}
		} else if strings.HasSuffix(valLower, "day") {
			if days, err := strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(valLower, "day"))); err == nil {
				return time.Duration(days) * 24 * time.Hour
			}
		} else if strings.HasSuffix(valLower, "d") {
			if days, err := strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(valLower, "d"))); err == nil {
				return time.Duration(days) * 24 * time.Hour
			}
		}

		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultValue
}

func GetConfig() *Config {
	if AppConfig == nil {
		AppConfig = loadConfig()
	}
	return AppConfig
}
