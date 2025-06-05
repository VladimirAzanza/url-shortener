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
	"github.com/VladimirAzanza/url-shortener/internal/repo"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLService struct {
	cfg        *config.Config
	storage    map[string]string
	repoSQLite repo.ISQLiteStorage
}

func NewURLService(cfg *config.Config, repoSQLite repo.ISQLiteStorage) *URLService {
	return &URLService{
		cfg:        cfg,
		storage:    make(map[string]string, 0),
		repoSQLite: repoSQLite,
	}
}

func (s *URLService) ShortenURL(ctx context.Context, originalURL string) string {
	shortID := generateUniqueID(originalURL)
	s.storage[shortID] = originalURL

	s.saveRecord(shortID, originalURL)
	return shortID
}

func (s *URLService) ShortenAPIURL(ctx context.Context, shortenRequest *dto.ShortenRequestDTO) string {
	return s.ShortenURL(ctx, shortenRequest.URL)
}

func (s *URLService) GetOriginalURL(ctx context.Context, shortID string) (string, bool) {
	timer := time.NewTimer(100 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-timer.C:
		originalURL, exists := s.storage[shortID]
		return originalURL, exists
	case <-ctx.Done():
		return "", false
	}
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

func (s *URLService) BatchShortenURL(ctx context.Context, request dto.BatchRequestDTO) (string, error) {
	shortID := s.ShortenURL(ctx, request.OriginalURL)
	err := s.repoSQLite.SaveBatchURL(ctx, shortID, request.OriginalURL)
	if err != nil {
		return "", fmt.Errorf("failed to save URL %s: %w", request.OriginalURL, err)
	}

	return shortID, nil
}

func generateUniqueID(originalURL string) string {
	hash := sha256.Sum256([]byte(originalURL))
	hashStr := hex.EncodeToString(hash[:])[:8]
	timestamp := fmt.Sprintf("%x", time.Now().UnixNano())[:8]
	return hashStr + timestamp
}
