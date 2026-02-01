package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dgentikian/rent-receipt-generator/internal/config"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func Connect(cfg *config.DatabaseConfig) (*DB, error) {
	// Build connection string with search_path
	password := cfg.Password
	if password == "" {
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=%s search_path=public",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Name,
			cfg.SSLMode,
		)
		return connectWithDSN(dsn, cfg)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=public",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)

	return connectWithDSN(dsn, cfg)
}

func connectWithDSN(dsn string, cfg *config.DatabaseConfig) (*DB, error) {
	// Debug: Log connection details (without password)
	fmt.Printf("Connecting to database: host=%s port=%d user=%s dbname=%s\n",
		cfg.Host, cfg.Port, cfg.User, cfg.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) Health() error {
	return db.Ping()
}
