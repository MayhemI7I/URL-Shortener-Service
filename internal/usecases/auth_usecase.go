package usecases

import (
	"context"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/services"
)

// AuthUseCase реализует интерфейс AuthService
type AuthUseCase struct {
	service services.AuthService
}


// NewAuthUseCase создает новый экземпляр сервиса аутентификации
func NewAuthUseCase(service services.AuthService) *AuthUseCase {
		return &AuthUseCase{service: service}
	}


// Register регистрирует нового пользователя
func (a *AuthUseCase) Register(ctx context.Context, refreshToken *models.RefreshToken)  error {
	return a.service.Register(ctx)
}

// RefreshTokens обновляет токен пользователя
func (a *AuthUseCase) RefreshToken(ctx context.Context, refreshToken *models.RefreshToken)  error {
	return a.service.RefreshToken(ctx, refreshToken)

}

