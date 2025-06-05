package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/VladimirAzanza/url-shortener/config"
	"github.com/VladimirAzanza/url-shortener/internal/dto"
	"github.com/VladimirAzanza/url-shortener/internal/repo"
)

type URLService struct {
	cfg           *config.Config
	memoryStorage map[string]string
	repoSQLite    repo.ISQLiteStorage
	repoFile      repo.IFileStorage
}

func NewURLService(cfg *config.Config, repoSQLite repo.ISQLiteStorage, repoFile repo.IFileStorage) *URLService {
	return &URLService{
		cfg:           cfg,
		memoryStorage: make(map[string]string, 0),
		repoSQLite:    repoSQLite,
		repoFile:      repoFile,
	}
}

func (s *URLService) ShortenURL(ctx context.Context, originalURL string) string {
	shortID := generateUniqueID(originalURL)
	s.memoryStorage[shortID] = originalURL

	s.repoFile.SaveRecord(shortID, originalURL)
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
		originalURL, exists := s.memoryStorage[shortID]
		return originalURL, exists
	case <-ctx.Done():
		return "", false
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
