package repo

import "context"

type ISQLiteStorage interface {
	SaveBatchURL(ctx context.Context, shortID, originalURL string) error
}
