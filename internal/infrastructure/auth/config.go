package auth

import (
	"fmt"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/config"
	"github.com/spf13/pflag"
)

// JWTConfig конфигурация для JWT
type JWTConfig struct {
	config.BaseConfig
	Secret                 string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
}

// NewJWTConfig создает новую конфигурацию JWT
func NewJWTConfig() *JWTConfig {
	return &JWTConfig{
		BaseConfig:            config.NewBaseConfig("JWT"),
		Secret:                "secret",
		AccessTokenExpiration:  24 * time.Hour,
		RefreshTokenExpiration: 720 * time.Hour, // 30 days
	}
}

// AddFlags добавляет флаги для конфигурации JWT
func (c *JWTConfig) AddFlags() {
	c.AddStringFlag(c.GetFlags(), "secret", c.Secret, "JWT secret")
	c.AddDurationFlag(c.GetFlags(), "token-expiration", c.AccessTokenExpiration, "JWT token expiration")
	c.AddDurationFlag(c.GetFlags(), "refresh-expiration", c.RefreshTokenExpiration, "JWT refresh token expiration")
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func (c *JWTConfig) LoadFromEnv() {
	// Базовая реализация уже обрабатывает переменные окружения
	// через методы Add*Flag
}

// Validate проверяет корректность конфигурации
func (c *JWTConfig) Validate() error {
	// Проверяем обязательные поля
	if err := c.ValidateRequired(map[string]interface{}{
		"secret": c.Secret,
	}); err != nil {
		return err
	}

	// Проверяем длительности
	if err := c.ValidateDuration(map[string]time.Duration{
		"token-expiration":   c.AccessTokenExpiration,
		"refresh-expiration": c.RefreshTokenExpiration,
	}); err != nil {
		return err
	}

	// Проверяем логические ограничения
	if c.RefreshTokenExpiration <= c.AccessTokenExpiration {
		return fmt.Errorf("время жизни refresh-токена должно быть больше времени жизни access-токена")
	}

	return nil
}

// GetDefaultConfig возвращает конфигурацию по умолчанию
func (c *JWTConfig) GetDefaultConfig() interface{} {
	return NewJWTConfig()
}

// GetJWTSecret возвращает секретный ключ JWT
func (c *JWTConfig) GetJWTSecret() string {
	return c.Secret
}

// GetAccessTokenExpiration возвращает время жизни access-токена
func (c *JWTConfig) GetAccessTokenExpiration() time.Duration {
	return c.AccessTokenExpiration
}

// GetRefreshTokenExpiration возвращает время жизни refresh-токена
func (c *JWTConfig) GetRefreshTokenExpiration() time.Duration {
	return c.RefreshTokenExpiration
}
