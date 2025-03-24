package config

import (
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/config"
)

// AuthConfig реализует интерфейс interfaces.AuthConfig
type AuthConfig struct {
	JWTSecret              string
	AccessTokenExpiration        time.Duration
	RefreshTokenExpiration time.Duration
}

// NewAuthConfig создает новый экземпляр AuthConfig
func NewAuthConfig(cfg *config.Config) interfaces.AuthConfig {
	return &AuthConfig{
		JWTSecret:              cfg.JWT.Secret,
		AccessTokenExpiration:        cfg.JWT.AccessTokenExpiration,
		RefreshTokenExpiration: cfg.JWT.RefreshTokenExpiration,
	}
}

// GetJWTSecret возвращает секретный ключ JWT
func (c *AuthConfig) GetJWTSecret() string {
	return c.JWTSecret
}

// GetTokenExpiration возвращает время жизни токена
func (c *AuthConfig) GetAccessTokenExpiration() time.Duration {
	return c.AccessTokenExpiration
}

// GetRefreshTokenExpiration возвращает время жизни Refresh токена 
func (c *AuthConfig) GetRefreshTokenExpiration() time.Duration {
		return c.RefreshTokenExpiration
}
