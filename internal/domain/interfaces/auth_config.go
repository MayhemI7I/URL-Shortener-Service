package interfaces

import "time"

// AuthConfig определяет интерфейс для конфигурации аутентификации
type AuthConfig interface {
	// GetJWTSecret возвращает секретный ключ для JWT токенов
	GetJWTSecret() string
	// GetTokenExpiration возвращает время жизни токена
	GetAccessTokenExpiration() time.Duration
	// GetRefreshTokenExpiration возвращает время жизни refresh токена
	GetRefreshTokenExpiration() time.Duration
	
}
