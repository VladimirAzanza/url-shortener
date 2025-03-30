package services

import (
	"crypto/sha256"
	"encoding/hex"
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
	shortID := generateUniqueID(originalURL)
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

func generateUniqueID(originalURL string) string {
	hash := sha256.Sum256([]byte(originalURL))
	hashStr := hex.EncodeToString(hash[:])[:8]
	timestamp := fmt.Sprintf("%x", time.Now().UnixNano())[:8]
	return hashStr + timestamp
}
