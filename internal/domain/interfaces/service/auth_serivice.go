package service

import (
	"context"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
)
// AuthService определяет бизнес-логику для аутентификации
type AuthService interface {
	// Register регистрирует нового пользователя
	Register(ctx context.Context)  (*models.User,error)

	// RefreshToken обновляет токен пользователя
	RefreshToken(ctx context.Context, user *models.RefreshToken) error

}