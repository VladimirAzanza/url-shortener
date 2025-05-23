package repo

import (
	"database/sql"
	"fmt"

	"github.com/VladimirAzanza/url-shortener/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func NewDB(cfg *config.Config) (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)

	switch cfg.DatabaseType {
	case "sqlite":
		db, err = sql.Open("sqlite3", cfg.DatabaseDSN)
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
