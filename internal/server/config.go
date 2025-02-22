package server

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// Allow missing .env file in production
		if os.Getenv("ENV") != "production" {
			return nil, nil
		}
	}

	config := &Config{
		Port: GetEnv("PORT", "3000"),
	}

	return config, nil
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func ParseDuration(duration string) time.Duration {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return 24 * time.Hour // Default to 24 hours
	}
	return d
} 