package memory

import (
	"github.com/VladimirAzanza/url-shortener/internal/repo"
)

type MemoryRepository struct {
	storage map[string]string
}

func NewMemoryRepository() repo.IMemoryStorage {
	return &MemoryRepository{
		storage: make(map[string]string, 0),
	}
}

func (r *MemoryRepository) SaveShortID(shortID, originalURL string) {
	r.storage[shortID] = originalURL
}

func (r *MemoryRepository) GetOriginalURL(shortID string) (string, bool) {
	originalURL, ok := r.storage[shortID]
	return originalURL, ok
}
