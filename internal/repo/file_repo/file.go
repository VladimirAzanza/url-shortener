package filerepo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/VladimirAzanza/url-shortener/config"
	"github.com/VladimirAzanza/url-shortener/internal/dto"
	"github.com/VladimirAzanza/url-shortener/internal/repo"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type FileRepository struct {
	cfg     *config.Config
	file    *os.File
	encoder *json.Encoder
	storage map[string]string
}

func NewFileRepository(cfg *config.Config) repo.IURLRepository {
	fileRepo := &FileRepository{
		cfg:     cfg,
		storage: make(map[string]string, 0),
	}
	fileRepo.initFile()
	return fileRepo
}

func (r *FileRepository) initFile() {
	file, err := os.OpenFile(r.cfg.FileStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open storage file")
		return
	}
	r.file = file
	r.encoder = json.NewEncoder(file)
}

func (r *FileRepository) SaveShortID(ctx context.Context, shortID, originalURL string) error {
	urlRecord := dto.URLRecord{
		UUID:        uuid.New().String(),
		ShortURL:    shortID,
		OriginalURL: originalURL,
	}
	r.storage[shortID] = originalURL

	if err := r.encoder.Encode(urlRecord); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}
	return nil
}

func (r *FileRepository) SaveBatchURL(ctx context.Context, shortID, originalURL string) error {
	return r.SaveShortID(ctx, shortID, originalURL)
}

// Implementation for searching in the file (may be inefficient for many records)
// In a real implementation, consider using a database or an index.
func (r *FileRepository) GetOriginalURL(ctx context.Context, shortID string) (string, bool, error) {
	originalURL, ok := r.storage[shortID]
	log.Info().Msg("In a real implementation, consider using a database or an index.")
	return originalURL, ok, nil
}

func (r *FileRepository) Ping(ctx context.Context) error {
	if r.file == nil {
		return fmt.Errorf("file storage not initialized")
	}
	return nil
}

func (r *FileRepository) GetShortIDByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	for shortID, url := range r.storage {
		if url == originalURL {
			return shortID, nil
		}
	}
	return "", nil
}
