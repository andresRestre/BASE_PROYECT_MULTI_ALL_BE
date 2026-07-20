package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration values loaded from environment variables.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret          string
	JWTExpirationHours string

	ServerPort string
}

// Load reads the .env file and returns a Config struct with all values.
func Load() *Config {
	// Try loading .env from multiple paths to support running from different directories
	// (e.g., from backend/ or backend/cmd/server/)
	_ = godotenv.Load()            // current directory
	_ = godotenv.Load("../../.env") // from cmd/server/ -> backend/.env

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "multicliente_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-key-change-me-in-production"),
		JWTExpirationHours: getEnv("JWT_EXPIRATION_HOURS", "24"),

		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
