package storage

import (
	"context"
	"github.com/MayhemI7I/URL-Shortener-Service/config"
	"github.com/MayhemI7I/URL-Shortener-Service/domain"
	// "github.com/MayhemI7I/URL-Shortener-Service/internal/storage/file"
	// "github.com/MayhemI7I/URL-Shortener-Service/internal/storage/memory"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/storage/postgres"
	"time"


)
 
type Storage interface {
	Get(ctx context.Context, shortUrl string, origURL string) (string, error)
	Save(ctx context.Context, shortUrl, longUrl, user_id string) error
	FindByOriginalURL(ctx context.Context, shortUrlstring, userId string) (string, error)

	GetUserAllURLs(ctx context.Context, userId string) ([]domain.URLData, error)
	GetUserIDByRefreshToken(ctx context.Context, refreshToken string)(string,error)
	GetNewAccessToken(ctx context.Context, refreshToken string)(string,string,error)
	MarkURLsAsDeleted(ctx context.Context, shortURLs []string, userId string) error
	
	SaveRefreshToken(ctx context.Context,refreshToken, UserID string, expiresAt time.Time)error
	DeleteRefreshToken(ctx context.Context,refreshToken string)error
	
	Close() error
}

func NewStorage(c config.Config) (Storage, error) {
	if c.DataBaseDSN != "" {
		return postgres.NewPostgresStorage(c.DataBaseDSN)
	}
	return nil, nil
	// if c.FileStorage != "" {
	// 	return file.NewFileStorage(c.FileStorage)
	// }

	// return memory.NewMemoryStorage()
}
