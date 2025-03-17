package services

import "fmt"

type URLService struct {
	storage map[string]string
}

func NewURLService() *URLService {
	return &URLService{
		storage: make(map[string]string),
	}
}

func (s *URLService) ShortenURL(originalURL string) string {
	shortID := "123"
	_, exists := s.storage[shortID]

	if exists {
		shortID += "10"
		s.storage[shortID] = originalURL
	} else {
		s.storage[shortID] = originalURL
	}

	return fmt.Sprintf("http://localhost:8080/%s", (shortID))
}
