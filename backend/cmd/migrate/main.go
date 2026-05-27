package main

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"
	"time"

	"your-next-game-backend/internal/platform/config"
	"your-next-game-backend/internal/platform/migrations"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	if cfg.PostgresDSN == "" {
		log.Fatal("POSTGRES_DSN is required")
	}

	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("postgres open failed: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("postgres ping failed: %v", err)
	}

	migrationsDir := filepath.Clean("migrations")
	if err := migrations.EnsureDirExists(migrationsDir); err != nil {
		log.Fatalf("invalid migrations dir: %v", err)
	}

	upCount, err := migrations.CountUpFiles(migrationsDir)
	if err != nil {
		log.Fatalf("failed to count migration files: %v", err)
	}
	log.Printf("found %d up migration files", upCount)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := migrations.ApplyUp(ctx, db, migrationsDir); err != nil {
		log.Fatalf("apply migrations failed: %v", err)
	}

	log.Printf("migrations applied successfully")
}
