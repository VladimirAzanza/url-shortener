package repo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/VladimirAzanza/url-shortener/config"
)

func NewDB(cfg *config.Config) (*sql.DB, error) {
	var db *sql.DB
	var err error

	switch cfg.DatabaseType {
	case "sqlite":
		db, err = sql.Open("sqlite", cfg.DatabaseDSN)
	case "postgres":
		db, err = sql.Open("postgres", cfg.DatabaseDSN)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DatabaseType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Connected to %s database", cfg.DatabaseType)
	return db, nil
}
