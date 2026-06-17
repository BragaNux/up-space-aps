package config

import (
	"os"
	"time"
)

// Config junta toda a configuracao da aplicacao lida das variaveis de ambiente
type Config struct {
	AppPort     string
	DatabaseURL string
	JWTSecret   string
	TokenTTL    time.Duration
	CORSOrigin  string
}

// Load le as variaveis de ambiente e aplica valores padrao pro ambiente de dev quando faltar alguma
func Load() *Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@db:5432/up_espaco?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "up-espaco-dev-secret-change-me"
	}

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "*"
	}

	return &Config{
		AppPort:     port,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
		TokenTTL:    7 * 24 * time.Hour,
		CORSOrigin:  corsOrigin,
	}
}
