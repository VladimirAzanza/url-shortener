package repo

import "context"

type ISQLiteStorage interface {
	SaveBatchURL(ctx context.Context, shortID, originalURL string) error
}

type IFileStorage interface {
	SaveRecord(shortID, originalURL string) error
}

type IMemoryStorage interface {
	SaveShortID(shortID, originalURL string)
	GetOriginalURL(shortID string) (string, bool)
}

type StorageType string

const (
	MemoryStorage StorageType = "memory"
	FileStorage   StorageType = "file"
	SQLiteStorage StorageType = "sqlite"
)
