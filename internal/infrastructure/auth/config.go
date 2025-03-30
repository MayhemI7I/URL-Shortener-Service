package auth

import (
	"fmt"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/config"
	"github.com/spf13/pflag"
)

// JWTConfig конфигурация для JWT
type JWTConfig struct {
	name                   string
	flags                  *pflag.FlagSet
	Secret                 string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
}

// NewJWTConfig создает новую конфигурацию JWT
func NewJWTConfig() *JWTConfig {
	return &JWTConfig{
		name:                   "JWT",
		flags:                  pflag.NewFlagSet("JWT", pflag.ExitOnError),
		Secret:                 "secret",
		AccessTokenExpiration:  24 * time.Hour,
		RefreshTokenExpiration: 720 * time.Hour, // 30 days
	}
}

// AddFlags добавляет флаги для конфигурации JWT
func (c *JWTConfig) AddFlags() {
	c.flags.StringVarP(&c.Secret, "jwt-secret", "", c.Secret, "JWT secret key")
	c.flags.DurationVarP(&c.AccessTokenExpiration, "jwt-token-expiration", "", c.AccessTokenExpiration, "JWT token expiration")
	c.flags.DurationVarP(&c.RefreshTokenExpiration, "jwt-refresh-expiration", "", c.RefreshTokenExpiration, "JWT refresh token expiration")
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func (c *JWTConfig) LoadFromEnv() {
	// Здесь можно добавить логику загрузки из env переменных
}

// Validate проверяет корректность конфигурации
func (c *JWTConfig) Validate() error {
	// Проверяем обязательные поля
	if c.Secret == "" {
		return fmt.Errorf("поле jwt-secret не может быть пустым")
	}

	// Проверяем длительности
	if c.AccessTokenExpiration <= 0 {
		return fmt.Errorf("поле jwt-token-expiration должно быть положительным числом")
	}
	if c.RefreshTokenExpiration <= 0 {
		return fmt.Errorf("поле jwt-refresh-expiration должно быть положительным числом")
	}

	// Проверяем логические ограничения
	if c.RefreshTokenExpiration <= c.AccessTokenExpiration {
		return fmt.Errorf("время жизни refresh-токена должно быть больше времени жизни access-токена")
	}

	return nil
}

// GetFlags возвращает набор флагов
func (c *JWTConfig) GetFlags() *pflag.FlagSet {
	return c.flags
}

// GetSecretKey возвращает секретный ключ JWT
func (c *JWTConfig) GetSecretKey() string {
	return c.Secret
}

// GetTokenExpiration возвращает время жизни токена в секундах
func (c *JWTConfig) GetTokenExpiration() int {
	return int(c.AccessTokenExpiration.Seconds())
}

// GetJWTSecret возвращает секретный ключ JWT (алиас для поддержки старого интерфейса)
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

// Проверка соответствия интерфейсу
var _ config.JWTConfigProvider = (*JWTConfig)(nil)
var _ config.Configurable = (*JWTConfig)(nil)
