package config

import "os"

// Config holds all application configuration.
type Config struct {
	DatabaseURL string
	Port        string
}

// Load reads configuration from environment variables.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/api_clientes?sslmode=disable"
	}

	return &Config{
		DatabaseURL: dbURL,
		Port:        port,
	}
}
