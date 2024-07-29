// File: internal/db/db.go
package db

import (
	"database/sql"
	"fmt"

	"github.com/JettZgg/LineUp/internal/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Initialize(cfg config.DatabaseConfig) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
