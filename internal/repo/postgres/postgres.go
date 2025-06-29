package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/VladimirAzanza/url-shortener/internal/repo"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostgreSQLRepository struct {
	db *sql.DB
}

func NewPostgreSQLRepository(db *sql.DB) repo.IURLRepository {
	return &PostgreSQLRepository{
		db: db,
	}
}

func (r *PostgreSQLRepository) SaveShortID(ctx context.Context, shortID, originalURL string) error {
	return r.SaveBatchURL(ctx, shortID, originalURL)
}

func (r *PostgreSQLRepository) SaveBatchURL(ctx context.Context, shortID, originalURL string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO short_urls (uuid, short_url, original_url) 
         VALUES ($1, $2, $3) 
         ON CONFLICT (original_url) DO NOTHING`,
		uuid.New().String(), shortID, originalURL)

	if err != nil {
		return fmt.Errorf("could not insert URL: %w", err)
	}

	return tx.Commit()
}

func (r *PostgreSQLRepository) GetShortIDByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	var shortID string
	err := r.db.QueryRowContext(ctx,
		"SELECT short_url FROM short_urls WHERE original_url = $1", originalURL).Scan(&shortID)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return shortID, nil
}

func (r *PostgreSQLRepository) GetOriginalURL(ctx context.Context, shortID string) (string, bool, error) {
	var originalURL string
	err := r.db.QueryRowContext(ctx,
		"SELECT original_url FROM short_urls WHERE short_url = $1",
		shortID).Scan(&originalURL)

	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	return originalURL, true, nil
}

func (r *PostgreSQLRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *PostgreSQLRepository) BatchDeleteURLs(ctx context.Context, shortURLs []string) error {
	query := `
        UPDATE short_urls 
        SET is_deleted = true 
        WHERE short_url = ANY($1) AND is_deleted = false`

	_, err := r.db.ExecContext(ctx, query, pq.Array(shortURLs))
	return err
}
