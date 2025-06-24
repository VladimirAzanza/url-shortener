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

func NewSQLiteRepository(db *sql.DB) repo.IURLRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) SaveShortID(ctx context.Context, shortID, originalURL string) error {
	return r.SaveBatchURL(ctx, shortID, originalURL)
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

func (r *SQLiteRepository) GetShortIDByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	var shortID string
	err := r.db.QueryRowContext(ctx,
		"SELECT short_url FROM short_urls WHERE original_url = ?", originalURL).Scan(&shortID)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return shortID, nil
}

func (r *SQLiteRepository) GetOriginalURL(ctx context.Context, shortID string) (string, bool, error) {
	var originalURL string
	err := r.db.QueryRowContext(ctx,
		"SELECT original_url FROM short_urls WHERE short_url = ?",
		shortID).Scan(&originalURL)

	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	return originalURL, true, nil
}

func (r *SQLiteRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *SQLiteRepository) BatchDeleteURLs(ctx context.Context, shortURLs []string) error {
	return nil
}
