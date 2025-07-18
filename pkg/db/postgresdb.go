package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	TripsTable     = "trips"
	RideFaresTable = "ride_fares"
	UsersTable     = "users"
	DriversTable   = "drivers"
)

// PostgresConfig holds PostgreSQL connection configuration
type PostgresConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	SSLMode  string
}

// NewPostgresDefaultConfig creates a new PostgreSQL configuration from environment variables
func NewPostgresDefaultConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     getEnvOrDefault("POSTGRES_HOST", "localhost"),
		Port:     getEnvOrDefault("POSTGRES_PORT", "5432"),
		Database: getEnvOrDefault("POSTGRES_DB", "bolt_app"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		SSLMode:  getEnvOrDefault("POSTGRES_SSLMODE", "disable"),
	}
}

// NewPostgresClient creates a new PostgreSQL client
func NewPostgresClient(ctx context.Context, cfg *PostgresConfig) (*sql.DB, error) {
	if cfg.Username == "" {
		return nil, fmt.Errorf("postgres username is required")
	}
	if cfg.Password == "" {
		return nil, fmt.Errorf("postgres password is required")
	}
	if cfg.Database == "" {
		return nil, fmt.Errorf("postgres database is required")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to PostgreSQL at %s:%s/%s", cfg.Host, cfg.Port, cfg.Database)
	return db, nil
}

// GetDB returns the database instance
func GetDB(db *sql.DB) *sql.DB {
	return db
}

// CloseDB safely closes the database connection
func CloseDB(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
