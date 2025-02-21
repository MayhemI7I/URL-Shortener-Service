package storage

import (
	"database/sql"
	"local/config"
)

type Storage interface {
	Get(ctx context.Context, shortUrl string) (string, error)
	Save(ctx context.Context, shortUrl, longUrl string) error
}

func NewStorage(c config.Config) (Storage, error) {
	if c.DataBaseDSN != "" {
		return NewPostgresStorage(c.DataBaseDSN)
	}
	if c.FileStorage != "" {
		return NewFileStorage(c.FileStorage)
	}

	return NewMemoryStorage()
}
	