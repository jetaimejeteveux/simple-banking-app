package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Host            string
	Port            string
	Debug           bool
	AllowOrigins    string
	ShutdownTimeout int

	// Add more configuration fields as needed
}

// NewConfig returns a new Config with default values
func NewConfig() (*Config, error) {
	return &Config{
		Host:            "0.0.0.0",
		Port:            "8080",
		Debug:           false,
		AllowOrigins:    "*",
		ShutdownTimeout: 30,
	}, nil
}

// LoadEnv loads environment variables from .env file
func LoadEnv() error {
	err := godotenv.Load()
	// It's fine if .env doesn't exist in production
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// GetEnv gets an environment variable with a fallback
func GetEnv(key string, fallback ...string) string {
	value := os.Getenv(key)
	if len(value) == 0 && len(fallback) > 0 {
		return fallback[0]
	}
	return value
}

// GetDatabaseDSN returns the DSN for connecting to the database
func GetDatabaseDSN() string {
	host := GetEnv("DB_HOST", "localhost")
	port := GetEnv("DB_PORT", "5432")
	user := GetEnv("DB_USER", "postgres")
	password := GetEnv("DB_PASSWORD", "postgres")
	dbname := GetEnv("DB_NAME", "gobackend")
	sslmode := GetEnv("DB_SSLMODE", "disable")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}
