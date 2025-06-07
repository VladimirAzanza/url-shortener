package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/VladimirAzanza/url-shortener/config"
	"github.com/VladimirAzanza/url-shortener/internal/dto"
	"github.com/VladimirAzanza/url-shortener/internal/repo"
	"github.com/rs/zerolog/log"
)

type URLService struct {
	cfg  *config.Config
	repo repo.IURLRepository
}

func NewURLService(cfg *config.Config, repo repo.IURLRepository) *URLService {
	return &URLService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *URLService) ShortenURL(ctx context.Context, originalURL string) (string, error) {
	existingShortID, err := s.repo.GetShortIDByOriginalURL(ctx, originalURL)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("error checking existing URL: %w", err)
	}

	if existingShortID != "" {
		return existingShortID, nil
	}

	shortID := generateUniqueID(originalURL)
	if err := s.repo.SaveShortID(ctx, shortID, originalURL); err != nil {
		fmt.Printf("Error saving URL: %v\n", err)
		return "", err
	}
	return shortID, nil
}

func (s *URLService) ShortenAPIURL(ctx context.Context, shortenRequest *dto.ShortenRequestDTO) (string, error) {
	return s.ShortenURL(ctx, shortenRequest.URL)
}

func (s *URLService) GetOriginalURL(ctx context.Context, shortID string) (string, bool) {
	timer := time.NewTimer(100 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-timer.C:
		originalURL, exists, err := s.repo.GetOriginalURL(ctx, shortID)
		if err != nil {
			log.Error().Err(err).Msg("Error getting original URL")
			return "", false
		}
		return originalURL, exists
	case <-ctx.Done():
		return "", false
	}
}

func (s *URLService) BatchShortenURL(ctx context.Context, request dto.BatchRequestDTO) (string, error) {
	existingShortID, err := s.repo.GetShortIDByOriginalURL(ctx, request.OriginalURL)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("error checking existing URL: %w", err)
	}

	if existingShortID != "" {
		return existingShortID, nil
	}

	shortID := generateUniqueID(request.OriginalURL)
	err = s.repo.SaveBatchURL(ctx, shortID, request.OriginalURL)
	if err != nil {
		return "", fmt.Errorf("failed to save URL %s: %w", request.OriginalURL, err)
	}

	return shortID, nil
}

func (s *URLService) PingDB(ctx context.Context) error {
	if s.repo == nil {
		return fmt.Errorf("no storage repository configured")
	}
	return s.repo.Ping(ctx)
}

func (s *URLService) GetStorageType() string {
	return s.cfg.StorageType
}

func generateUniqueID(originalURL string) string {
	hash := sha256.Sum256([]byte(originalURL))
	hashStr := hex.EncodeToString(hash[:])[:8]
	timestamp := fmt.Sprintf("%x", time.Now().UnixNano())[:8]
	return hashStr + timestamp
}
