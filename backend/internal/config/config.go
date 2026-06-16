package config

import (
	"os"
)

type Config struct {
	AppPort     string
	DatabaseURL string
}

func Load() *Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@db:5432/up_espaco?sslmode=disable"
	}

	return &Config{
		AppPort:     port,
		DatabaseURL: databaseURL,
	}
}
