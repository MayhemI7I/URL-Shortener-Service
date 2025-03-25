package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// DB представляет общее подключение к базе данных
type DB struct {
	*sqlx.DB
}

// NewDB создаёт новое подключение к PostgreSQL
func NewDB(dsn string) (*DB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		logger.Log.Error("failed to open database", zap.Error(err))
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logger.Log.Error("failed to ping database", zap.Error(err))
		return nil, err
	}

	// Инициализация таблиц
	tables := []string{
		`CREATE TABLE IF NOT EXISTS short_urls (
			id SERIAL PRIMARY KEY,
			user_id UUID NOT NULL,
			short_url VARCHAR(255) UNIQUE NOT NULL,
			original_url VARCHAR(255) NOT NULL,
			deleted_flag BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id SERIAL PRIMARY KEY,
			user_id UUID NOT NULL,
			refresh_token VARCHAR(255) UNIQUE NOT NULL,
			expires_at TIMESTAMP NOT NULL
		)`,
	}

	for _, query := range tables {
		if _, err := db.Exec(query); err != nil {
			logger.Log.Error("failed to create table", zap.Error(err))
			return nil, err
		}
	}

	// Создание индексов
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_short_urls_user_id ON short_urls(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at)`,
	}
	for _, query := range indexes {
		if _, err := db.Exec(query); err != nil {
			logger.Log.Error("failed to create index", zap.Error(err))
			return nil, err
		}
	}

	logger.Log.Info("database, tables, and indexes initialized successfully")
	return &DB{db}, nil
}

// Close закрывает соединение с базой данных
func (db *DB) Close() error {
	err := db.DB.Close()
	if err != nil {
		logger.Log.Error("failed to close database connection", zap.Error(err))
	}
	return err
}

type PostgresStorage struct {
	db *DB
}

// NewPostgresURLStorage создаёт новое хранилище URL
func NewPostgresURLStorage(db *DB) interfaces.URLRepository {
	return &PostgresStorage{db: db}
}

// Get возвращает длинный URL по короткому
func (pg *PostgresStorage) Get(ctx context.Context, shortURL string, userID string) (string, error) {
	var origURL string
	query := `SELECT original_url FROM short_urls WHERE short_url = $1 AND user_id = $2`
	err := pg.db.GetContext(ctx, &origURL, query, shortURL, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Debug("short URL not found", zap.String("short_url", shortURL), zap.String("user_id", userID))
			return "", domain.ErrURLNotFound
		}
		logger.Log.Error("failed to get original URL", zap.Error(err))
		return "", err
	}
	return origURL, nil
}

// ListByUser возвращает все URL пользователя
func (pg *PostgresStorage) ListByUser(ctx context.Context, userID string) ([]*models.URLData, error) {
	var urls []*models.URLData
	query := `SELECT user_id, short_url, original_url, created_at FROM short_urls WHERE user_id = $1`
	err := pg.db.SelectContext(ctx, &urls, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Debug("no URLs found for user", zap.String("user_id", userID), zap.Error(err))
			return nil, nil
		}
		logger.Log.Error("failed to get user URLs", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	logger.Log.Debug("retrieved URLs for user", zap.String("user_id", userID), zap.Int("count", len(urls)))
	return urls, nil
}

// AddUserURL сохраняет пару короткий и длинный URL пользователя
func (pg *PostgresStorage) AddUserURL(ctx context.Context, url *models.URLData) error {
	query := `
		INSERT INTO short_urls (short_url, original_url, user_id) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (short_url) 
		DO UPDATE SET original_url = EXCLUDED.original_url 
		WHERE short_urls.user_id = $3
	`
	result, err := pg.db.ExecContext(ctx, query, url.ShortURL, url.OrigURL, url.User.ID)
	if err != nil {
		logger.Log.Error("failed to save short URL", zap.String("short_url", url.ShortURL), zap.Error(err))
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrURLExists // Если конфликт и ничего не обновлено
	}
	logger.Log.Debug("short URL saved", zap.String("short_url", url.ShortURL), zap.String("user_id", url.User.ID))
	return nil
}

// FindByOriginalURL ищет короткий URL по длинному
func (pg *PostgresStorage) FindByOriginalURL(ctx context.Context, originalURL string) (*models.URLData, error) {
	var data models.URLData
	query := `SELECT user_id, short_url, original_url, created_at FROM short_urls WHERE original_url = $1`
	err := pg.db.GetContext(ctx, &data, query, originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrURLNotFound
		}
		logger.Log.Error("failed to find short URL by original URL", zap.String("original_url", originalURL), zap.Error(err))
		return nil, err
	}
	return &data, nil
}

// FindByShortURL ищет длинный URL по короткому
func (pg *PostgresStorage) FindByShortURL(ctx context.Context, shortURL string) (*models.URLData, error) {
	var data models.URLData
	query := `SELECT user_id, short_url, original_url, created_at FROM short_urls WHERE short_url = $1`
	err := pg.db.GetContext(ctx, &data, query, shortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrURLNotFound
		}
		logger.Log.Error("failed to find short URL by original URL", zap.String("original_url", data.OrigURL), zap.Error(err))
		return nil, err
	}
	return &data, nil
}

// MarkAsDeleted помечает URL как удаленный
func (pg *PostgresStorage) MarkAsDeleted(ctx context.Context, shortURL string) error {
	query := `UPDATE short_urls SET deleted_flag = TRUE WHERE short_url = $1`
	result, err := pg.db.ExecContext(ctx, query, shortURL)
	if err != nil {
		logger.Log.Error("failed to mark URL as deleted", zap.String("short_url", shortURL), zap.Error(err))
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrURLNotFound
	}
	logger.Log.Debug("URL marked as deleted", zap.String("short_url", shortURL))
	return nil
}
