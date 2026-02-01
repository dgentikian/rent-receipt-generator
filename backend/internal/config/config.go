package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	App      AppConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
	Expiry string
}

type AppConfig struct {
	UploadsDir  string
	FrontendURL string
	BackendURL  string
	CORSOrigins string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	// Try .env.local first (for local development), then .env
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load() // Fallback to .env

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "rent_user"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "rent_receipts"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
			Expiry: getEnv("JWT_EXPIRY", "24h"),
		},
		App: AppConfig{
			UploadsDir:  getEnv("UPLOADS_DIR", "./uploads"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
			BackendURL:  getEnv("BACKEND_URL", "http://localhost:8080"),
			CORSOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		},
	}

	// Validate required fields
	// DB_PASSWORD is only required in production
	if config.Server.Env == "production" && config.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required in production")
	}
	if config.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
