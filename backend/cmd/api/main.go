package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"up-espaco/backend/internal/adapters/database/postgres"
	"up-espaco/backend/internal/adapters/http/middleware"
	"up-espaco/backend/internal/adapters/http/routes"
	"up-espaco/backend/internal/adapters/repositories"
	"up-espaco/backend/internal/config"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	store := postgres.NewDB(db)
	repos := repositories.NewRepositoryContainer(store)
	handler := routes.NewRouter(repos)

	wrapped := middleware.WithMetrics(middleware.WithCORS(handler))

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      wrapped,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("starting server on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
