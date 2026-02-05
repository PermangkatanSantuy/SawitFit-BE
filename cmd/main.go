package main

import (
	"context"
	"log"
	"os"

	"github.com/cmd/api"
	"github.com/internal/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("ParseConfig error:", err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close()

	// Example query to test connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)

	jwks := os.Getenv("SUPABASE_JWKS_URL")
	verifier, err := auth.NewVerifier(jwks)
	if err != nil {
		log.Fatal(err)
	}

	// Jalankan
	server := api.NewAPIServer(":8080", conn, verifier)
	if err := server.Run(); err != nil {
		log.Fatal("Server error:", err)
	}
}
