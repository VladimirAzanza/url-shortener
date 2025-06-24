package memory

import (
	"context"

	"github.com/VladimirAzanza/url-shortener/internal/repo"
	"github.com/rs/zerolog/log"
)

type MemoryRepository struct {
	storage map[string]string
}

func NewMemoryRepository() repo.IURLRepository {
	return &MemoryRepository{
		storage: make(map[string]string, 0),
	}
}

func (r *MemoryRepository) SaveShortID(ctx context.Context, shortID, originalURL string) error {
	r.storage[shortID] = originalURL
	return nil
}

func (r *MemoryRepository) SaveBatchURL(ctx context.Context, shortID, originalURL string) error {
	return r.SaveShortID(ctx, shortID, originalURL)
}

func (r *MemoryRepository) GetOriginalURL(ctx context.Context, shortID string) (string, bool, error) {
	originalURL, ok := r.storage[shortID]
	return originalURL, ok, nil
}

func (r *MemoryRepository) Ping(ctx context.Context) error {
	log.Info().Msg("Memory Storage always available")
	return nil
}

func (r *MemoryRepository) GetShortIDByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	for shortID, url := range r.storage {
		if url == originalURL {
			return shortID, nil
		}
	}
	return "", nil
}

func (r *MemoryRepository) BatchDeleteURLs(ctx context.Context, shortURLs []string) error {
	return nil
}
