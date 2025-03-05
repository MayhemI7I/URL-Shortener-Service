package storage

import (
	"context"
	"local/config"
	"local/domain"
	"local/internal/storage/file"
	"local/internal/storage/memory"
	"local/internal/storage/postgres"
	"time"


)
 
type Storage interface {
	Get(ctx context.Context, shortUrl string, origURL string) (string, error)
	Save(ctx context.Context, shortUrl, longUrl, user_id string) error
	FindByLongURL(ctx context.Context, shortUrlstring, userId string) (string, error)
	GetUserURLs(ctx context.Context, userId string) ([]domain.URLData, error)
	GetUserIDByRefreshToken(ctx context.Context, refreshToken string)(string,error)
	SaveRefreshToken(ctx context.Context,refreshToken, UserID string, expiresAt time.Time)error
	DeleteRefreshToken(ctx context.Context,refreshToken string)error
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
