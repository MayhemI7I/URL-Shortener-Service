package services

import (
	"context"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/models"
	"github.com/google/uuid"

	"github.com/MayhemI7I/URL-Shortener-Service/pkg/utils/jwtutil"
)

// AuthService реализует интерфейс interfaces.AuthService
type AuthService struct {
	repo interfaces.UserRepository
	cfg  interfaces.AuthConfig
}

// NewAuthService создает новый экземпляр AuthService
func NewAuthService(repo interfaces.UserRepository, cfg interfaces.AuthConfig) interfaces.AuthService {
	return &AuthService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *AuthService) Register(ctx context.Context) (*models.User,error) {
	user := &models.User{}

	id := uuid.New().String()

	refreshTokenStr, err := jwtutil.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	user.ID = id
	user.CreatedAt = time.Now().UTC()

	// Создаем объект RefreshToken
	refreshToken := &models.RefreshToken{}
	refreshToken.ID = id
	refreshToken.Token = refreshTokenStr
	refreshToken.ExpiresAt = time.Now().Add(s.cfg.GetRefreshTokenExpiration())
	return user,s.repo.SaveRefreshToken(ctx, refreshToken)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error {
	err := s.repo.SaveRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	return nil
}
