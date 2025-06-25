package services

import (
	"context"

	"github.com/VladimirAzanza/url-shortener/internal/dto"
)

type IURLService interface {
	ConcurrentBatchDelete(ctx context.Context, shortURLs []string) error
	BatchDeleteURLs(ctx context.Context, shortURLs []string) error
	ShortenURL(ctx context.Context, originalURL string) (string, error)
	ShortenAPIURL(ctx context.Context, shortenRequest *dto.ShortenRequestDTO) (string, error)
	GetOriginalURL(ctx context.Context, shortID string) (string, bool)
	BatchShortenURL(ctx context.Context, request dto.BatchRequestDTO) (string, error)
	PingDB(ctx context.Context) error
	GetStorageType() string
}
