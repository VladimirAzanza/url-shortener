package services

import (
	"fmt"
	"time"
)

type URLService struct {
	storage map[string]string
}

func NewURLService() *URLService {
	return &URLService{
		storage: make(map[string]string, 0),
	}
}

func (s *URLService) ShortenURL(baseURL string, originalURL string) string {
	shortID := generateUniqueId(originalURL)
	s.storage[shortID] = originalURL

	return fmt.Sprintf(baseURL, shortID)
}

func (s *URLService) GetOriginalURL(shortID string) (string, bool) {
	originalURL, exists := s.storage[shortID]
	if !exists {
		return "", false
	}

	return originalURL, true
}

func generateUniqueId(originalURL string) string {
	return fmt.Sprintf("%x%x", originalURL, time.Now().UnixNano())[:16]
}
