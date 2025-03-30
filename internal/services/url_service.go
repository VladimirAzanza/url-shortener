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

func (s *URLService) ShortenURL(originalURL string) string {
	shortID := generateUniqueId(originalURL)
	s.storage[shortID] = originalURL

	return shortID
}

func (s *URLService) GetOriginalURL(shortID string) (string, bool) {
	originalURL, exists := s.storage[shortID]
	if !exists {
		return "", false
	}

	return originalURL, true
}

func generateUniqueId(originalURL string) string {
	f := fmt.Sprintf("%x", originalURL)[:8]
	s := fmt.Sprintf("%x", time.Now().UnixNano())[:8]
	return f + s
}
