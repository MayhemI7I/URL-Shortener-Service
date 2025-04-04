package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/infrastructure"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/repository"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"
	"github.com/MayhemI7I/URL-Shortener-Service/pkg/utils/jwtutil"

	"go.uber.org/zap"
)

type PostgresUserStorage struct {
	db  *DB
	cfg infrastructure.JWTConfigProvider
}

// NewPostgresUserStorage создаёт новое хранилище пользователей
func NewPostgresUserStorage(db *DB, cfg infrastructure.JWTConfigProvider) repository.UserRepository {
	return &PostgresUserStorage{db: db, cfg: cfg}
}

// GetUserById возвращает пользователя по ID
func (pg *PostgresUserStorage) GetUserById(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	query := `SELECT id FROM users WHERE id = $1`
	err := pg.db.GetContext(ctx, &user.ID, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		logger.Log.Error("failed to get user by ID", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	return &user, nil
}

// GetUserByRefreshToken возвращает пользователя по refresh-токену
func (pg *PostgresUserStorage) GetUserByRefreshToken(ctx context.Context, refreshToken string) (*models.User, error) {
	var user models.User
	query := `SELECT user_id FROM refresh_tokens WHERE refresh_token = $1`
	err := pg.db.GetContext(ctx, &user.ID, query, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrTokenNotFound
		}
		logger.Log.Error("failed to get user by refresh token", zap.String("refresh_token", refreshToken), zap.Error(err))
		return nil, err
	}
	return &user, nil
}

// GetNewAccessToken retrieves a new access token and refresh token for the given refresh token.
// If the refresh token is expired, it generates a new refresh token and saves it to the database.
func (pg *PostgresUserStorage) GetNewAccessToken(ctx context.Context, refreshToken string) (string, string, error) {
	var (
		userID    string
		expiresAt time.Time
	)
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE refresh_token = $1`
	// Используем указатели, чтобы данные записывались в переменные
	err := pg.db.QueryRowContext(ctx, query, refreshToken).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", models.ErrTokenNotFound
		}
		logger.Log.Error("ошибка при запросе refresh-токена из базы", zap.Error(err))
		return "", "", err
	}
	newRefreshToken := refreshToken

	if time.Now().After(expiresAt) {
		logger.Log.Debug("refresh token expired", zap.String("refresh_token", refreshToken))
		newRefreshToken, err = jwtutil.GenerateRefreshToken()
		if err != nil {
			return "", "", err
		}
		newExpiresAt := time.Now().Add(jwtutil.RefreshTokenExpiration)

		// Создаем объект RefreshToken для сохранения
		tokenObj := &models.RefreshToken{
			User:      models.User{ID: userID},
			Token:     newRefreshToken,
			ExpiresAt: newExpiresAt,
		}

		err = pg.SaveRefreshToken(ctx, tokenObj)
		if err != nil {
			logger.Log.Error(err)
			return "", "", err
		}
	}

	accessToken, err := jwtutil.GenerateAccessToken(userID, pg.cfg.GetSecretKey())
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// SaveRefreshToken сохраняет refresh-токен
func (pg *PostgresUserStorage) SaveRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (refresh_token, user_id, expires_at) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (refresh_token) 
		DO UPDATE SET expires_at = EXCLUDED.expires_at, user_id = EXCLUDED.user_id
	`
	_, err := pg.db.ExecContext(ctx, query, refreshToken.Token, refreshToken.User.ID, refreshToken.ExpiresAt)
	if err != nil {
		logger.Log.Error("failed to save refresh token", zap.String("refresh_token", refreshToken.Token), zap.Error(err))
		return err
	}
	logger.Log.Debug("refresh token saved", zap.String("refresh_token", refreshToken.Token), zap.String("user_id", refreshToken.User.ID))
	return nil
}

// DeleteRefreshToken удаляет refresh-токен
func (pg *PostgresUserStorage) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	query := `DELETE FROM refresh_tokens WHERE refresh_token = $1`
	result, err := pg.db.ExecContext(ctx, query, refreshToken)
	if err != nil {
		logger.Log.Error("failed to delete refresh token", zap.String("refresh_token", refreshToken), zap.Error(err))
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.ErrTokenNotFound
	}
	logger.Log.Debug("refresh token deleted", zap.String("refresh_token", refreshToken))
	return nil
}
