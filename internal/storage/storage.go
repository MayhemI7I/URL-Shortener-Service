package storage

import (
	"local/config"
	"context"
	"local/internal/storage/postgres"
	"local/internal/storage/memory"
	"local/internal/storage/file"
)

type Storage interface {
	Get(ctx context.Context, shortUrl string) (string, error)
	Save(ctx context.Context, shortUrl, longUrl string) error
	Close()error
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
