package postgres

import (
	"context"
	"database/sql"
	"time"
   "os"
   "errors"

	"local/domain"
	"local/utils/jwtutil"
	"local/logger"
   
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	db *sqlx.DB
}

var secretKey = (os.Getenv("JWT-SECRET")) // Лучше вынести в конфигурацию

// NewPostgresStorage создаёт новое подключение к PostgreSQL
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
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
			long_url VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id SERIAL PRIMARY KEY,
			refresh_token VARCHAR(255) UNIQUE NOT NULL,
			user_id UUID NOT NULL,
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
	return &PostgresStorage{db: db}, nil
}

// Close закрывает соединение с базой данных
func (pg *PostgresStorage) Close() error {
	err := pg.db.Close()
	if err != nil {
		logger.Log.Error("failed to close database connection", zap.Error(err))
	}
	return err
}

// Get возвращает длинный URL по короткому
func (pg *PostgresStorage) Get(ctx context.Context, shortURL string, userID string) (string, error) {
	var longURL string
	query := `SELECT long_url FROM short_urls WHERE short_url = $1 AND user_id = $2`
	err := pg.db.GetContext(ctx, &longURL, query, shortURL, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Debug("short URL not found", zap.String("short_url", shortURL), zap.String("user_id", userID))
			return "", domain.ErrURLNotFound
		}
		logger.Log.Error("failed to get long URL", zap.Error(err))
		return "", err
	}
	return longURL, nil
}

// GetUserURLs возвращает все URL пользователя
func (pg *PostgresStorage) GetUserAllURLs(ctx context.Context, userID string) ([]domain.URLData, error) {
	var urls []domain.URLData
	query := `SELECT short_url, long_url FROM short_urls WHERE user_id = $1`
	err := pg.db.SelectContext(ctx, &urls, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Пустой список, а не ошибка
		}
		logger.Log.Error("failed to get user URLs", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	logger.Log.Debug("retrieved URLs for user", zap.String("user_id", userID), zap.Int("count", len(urls)))
	return urls, nil
}

// Save сохраняет короткий URL
func (pg *PostgresStorage) Save(ctx context.Context, shortURL, origURL, userID string) error {
	query := `
		INSERT INTO short_urls (short_url, long_url, user_id) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (short_url) 
		DO UPDATE SET long_url = EXCLUDED.long_url 
		WHERE short_urls.user_id = $3
	`
	result, err := pg.db.ExecContext(ctx, query, shortURL, origURL, userID)
	if err != nil {
		logger.Log.Error("failed to save short URL", zap.String("short_url", shortURL), zap.Error(err))
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrURLExists // Если конфликт и ничего не обновлено
	}
	logger.Log.Debug("short URL saved", zap.String("short_url", shortURL), zap.String("user_id", userID))
	return nil
}

// FindByLongURL ищет короткий URL по длинному
func (pg *PostgresStorage) FindByLongURL(ctx context.Context, longURL, userID string) (string, error) {
	var shortURL string
	query := `SELECT short_url FROM short_urls WHERE long_url = $1 AND user_id = $2`
	err := pg.db.GetContext(ctx, &shortURL, query, longURL, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrURLNotFound
		}
		logger.Log.Error("failed to find short URL by long URL", zap.String("long_url", longURL), zap.String("user_id", userID), zap.Error(err))
		return "", err
	}
	return shortURL, nil
}

// GetUserIDByRefreshToken возвращает user_id по refresh-токену
func (pg *PostgresStorage) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	var userID string
	query := `SELECT user_id FROM refresh_tokens WHERE refresh_token = $1`
	err := pg.db.GetContext(ctx, &userID, query, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrTokenNotFound
		}
		logger.Log.Error("failed to get user ID by refresh token", zap.String("refresh_token", refreshToken), zap.Error(err))
		return "", err
	}
	return userID, nil
}

// GetNewAccessToken генерирует новый access-токен
func (pg *PostgresStorage) GetNewAccessToken(ctx context.Context, refreshToken string) (string, string, error) {
	var (
		userID    string
		expiresAt time.Time
	)
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE refresh_token = $1`
	err := pg.db.GetContext(ctx, &struct {
		UserID    string    `db:"user_id"`
		ExpiresAt time.Time `db:"expires_at"`
	}{UserID: userID, ExpiresAt: expiresAt}, query, refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", domain.ErrTokenNotFound
		}
		return "", "", err
	}

	if time.Now().After(expiresAt) {
		return "", "", domain.ErrTokenExpired
	}

	accessToken, err := jwtutil.GenerateAccessToken(userID, "secret") // Здесь используется jwtutil
	if err != nil {
		return "", "", err
	}

	// Здесь можно сгенерировать новый refresh-токен, если требуется
	newRefreshToken := refreshToken // Пока используем старый
	return accessToken, newRefreshToken, nil
}

// SaveRefreshToken сохраняет refresh-токен
func (pg *PostgresStorage) SaveRefreshToken(ctx context.Context, refreshToken, userID string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (refresh_token, user_id, expires_at) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (refresh_token) 
		DO UPDATE SET expires_at = EXCLUDED.expires_at, user_id = EXCLUDED.user_id
	`
	_, err := pg.db.ExecContext(ctx, query, refreshToken, userID, expiresAt)
	if err != nil {
		logger.Log.Error("failed to save refresh token", zap.String("refresh_token", refreshToken), zap.Error(err))
		return err
	}
	logger.Log.Debug("refresh token saved", zap.String("refresh_token", refreshToken), zap.String("user_id", userID))
	return nil
}

// DeleteRefreshToken удаляет refresh-токен
func (pg *PostgresStorage) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	query := `DELETE FROM refresh_tokens WHERE refresh_token = $1`
	result, err := pg.db.ExecContext(ctx, query, refreshToken)
	if err != nil {
		logger.Log.Error("failed to delete refresh token", zap.String("refresh_token", refreshToken), zap.Error(err))
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrTokenNotFound
	}
	logger.Log.Debug("refresh token deleted", zap.String("refresh_token", refreshToken))
	return nil
}

