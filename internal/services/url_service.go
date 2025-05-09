package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/VladimirAzanza/url-shortener/config"
	"github.com/VladimirAzanza/url-shortener/internal/dto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLService struct {
	cfg     *config.Config
	storage map[string]string
}

func NewURLService(cfg *config.Config) *URLService {
	return &URLService{
		cfg:     cfg,
		storage: make(map[string]string, 0),
	}
}

func (s *URLService) ShortenURL(originalURL string) string {
	shortID := generateUniqueID(originalURL)
	s.storage[shortID] = originalURL

	s.saveRecord(shortID, originalURL)

	return shortID
}

func (s *URLService) ShortenAPIURL(ctx context.Context, shortenRequest *dto.ShortenRequestDTO) string {
	return s.ShortenURL(shortenRequest.URL)
}

func (s *URLService) GetOriginalURL(shortID string) (string, bool) {
	originalURL, exists := s.storage[shortID]
	return originalURL, exists
}

func (s *URLService) saveRecord(shortID, originalURL string) {
	if s.cfg.FileStoragePath == "" {
		log.Error().Msg("File storage path is empty")
		return
	}

	urlRecord := URLRecord{
		UUID:        uuid.New().String(),
		ShortURL:    shortID,
		OriginalURL: originalURL,
	}

	file, err := os.OpenFile(s.cfg.FileStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open storage file")
		return
	}
	defer file.Close()

	data, err := json.Marshal(urlRecord)
	if err != nil {
		log.Error().Err(err).Msg("Failed to Marshal url record")
		return
	}

	if _, err := file.WriteString(string(data) + "\n"); err != nil {
		log.Error().Err(err).Msg("Failed to write record")
		return
	}
}

func generateUniqueID(originalURL string) string {
	hash := sha256.Sum256([]byte(originalURL))
	hashStr := hex.EncodeToString(hash[:])[:8]
	timestamp := fmt.Sprintf("%x", time.Now().UnixNano())[:8]
	return hashStr + timestamp
}
