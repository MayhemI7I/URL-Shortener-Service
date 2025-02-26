package storage

import (
	"context"
	"local/config"
	"local/internal/storage/file"
	"local/internal/storage/memory"
	"local/internal/storage/postgres"
)

type Storage interface {
	Get(ctx context.Context, shortUrl string) (string, error)
	Save(ctx context.Context, shortUrl, longUrl string) error
	FindByLongURL(context.Context, string) (string, error)
	Close() error
}

func NewStorage(c config.Config) (Storage, error) {
	if c.DataBaseDSN != "" {
		return postgres.NewPostgresStorage(c.DataBaseDSN)
	}
	if c.FileStorage != "" {
		return file.NewFileStorage(c.FileStorage)
	}

	return memory.NewMemoryStorage()
}
