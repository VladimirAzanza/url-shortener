package services

import (
	"testing"
)

func TestShortenURL(t *testing.T) {
	tests := []struct {
		name        string
		originalURL string
		want        string
	}{
		{"Valid URL", "https://myurl_example.com", "http://localhost:8080/123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewURLService()
			got := s.ShortenURL(tt.originalURL)
			if got != tt.want {
				t.Errorf("ShortenURL() = %v, expected empty: %v", got, tt.want)
			}
		})
	}
}

func TestGetOriginalURL(t *testing.T) {
	tests := []struct {
		name    string
		shortID string
		want    string
	}{
		{"Existing URL", "123", "https://myurl_example.com"},
		{"Non-existing URL", "999", ""},
	}

	s := NewURLService()
	s.ShortenURL("https://myurl_example.com")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalURL, found := s.GetOriginalURL(tt.shortID)
			if found {
				if originalURL != "https://myurl_example.com" {
					t.Errorf("GetOriginalURL() found = %v, want %v", originalURL, "https://myurl_example.com")
				}
			} else if tt.want != "" {
				t.Errorf("GetOriginalURL() did not find URL, expected %v", tt.want)
			}
		})
	}
}
