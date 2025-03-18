package services

import "testing"

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
