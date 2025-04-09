package services

import (
	"context"
	"testing"

	"github.com/VladimirAzanza/url-shortener/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL(t *testing.T) {
	tests := []struct {
		name        string
		originalURL string
	}{
		{"Valid URL", "https://example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewURLService()
			req := &dto.ShortenRequestDTO{URL: tt.originalURL}
			shortID := s.ShortenURL(context.Background(), req)

			assert.NotEmpty(t, shortID)
			assert.Len(t, shortID, 16)

			originalURL, exists := s.GetOriginalURL(shortID)
			assert.True(t, exists)
			assert.Equal(t, tt.originalURL, originalURL)
		})
	}
}

func TestGetOriginalURL(t *testing.T) {
	s := NewURLService()
	testURL := "https://example.com"
	req := &dto.ShortenRequestDTO{URL: testURL}
	shortID := s.ShortenURL(context.Background(), req)

	tests := []struct {
		name     string
		shortID  string
		expected string
		exists   bool
	}{
		{"Existing URL", shortID, testURL, true},
		{"Non Existin URL", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalURL, exists := s.GetOriginalURL(tt.shortID)
			assert.Equal(t, tt.exists, exists)
			if exists {
				assert.Equal(t, tt.expected, originalURL)
			}
		})
	}
}

func TestGenerateUniqueId(t *testing.T) {
	fURL := "https://example.com"
	sURL := "https://google.com"

	id1 := generateUniqueID(fURL)
	id2 := generateUniqueID(sURL)

	assert.Len(t, id1, 16)
	assert.Len(t, id2, 16)
	assert.NotEqual(t, id1, id2, "ID's should be unique for different URLS")
}
