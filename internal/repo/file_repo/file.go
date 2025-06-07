package filerepo

import (
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
	cfg *config.Config
}

func NewFileRepository(cfg *config.Config) repo.URLRepository {
	return &FileRepository{
		cfg: cfg,
	}
}

func (s *FileRepository) SaveRecord(shortID, originalURL string) error {
	if s.cfg.FileStoragePath == "" {
		log.Error().Msg("File storage path is empty")
		return fmt.Errorf("file storage path is empty")
	}

	urlRecord := dto.URLRecord{
		UUID:        uuid.New().String(),
		ShortURL:    shortID,
		OriginalURL: originalURL,
	}

	file, err := os.OpenFile(s.cfg.FileStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open storage file")
		return fmt.Errorf("failed to open storage file")
	}
	defer file.Close()

	data, err := json.Marshal(urlRecord)
	if err != nil {
		log.Error().Err(err).Msg("Failed to Marshal url record")
		return fmt.Errorf("failed to Marshal url record")
	}

	if _, err := file.WriteString(string(data) + "\n"); err != nil {
		log.Error().Err(err).Msg("Failed to write record")
		return fmt.Errorf("failed to write record")
	}
	return nil
}
