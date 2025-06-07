package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/VladimirAzanza/url-shortener/internal/repo"
	"github.com/google/uuid"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) repo.URLRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) SaveBatchURL(ctx context.Context, shortID, originalURL string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO short_urls (uuid, short_url, original_url) 
         VALUES (?, ?, ?)`,
		uuid.New().String(), shortID, originalURL)

	if err != nil {
		return fmt.Errorf("could not insert URL: %w", err)
	}

	return tx.Commit()
}
