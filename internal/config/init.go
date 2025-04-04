package config

import (
	"fmt"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/infrastructure"
)

// Init создает и инициализирует конфигурацию приложения с пользовательскими провайдерами
func Init(
	httpConfig infrastructure.HTTPConfigProvider,
	loggerConfig infrastructure.LoggerConfigProvider,
	dbConfig infrastructure.DBConfigProvider,
	fileStorageConfig infrastructure.FileStorageConfigProvider,
	jwtConfig infrastructure.JWTConfigProvider,
) (*AppConfig, error) {
	// Создаем конфигурацию с пользовательскими провайдерами
	cfg := NewAppConfig(
		httpConfig,
		loggerConfig,
		dbConfig,
		fileStorageConfig,
		jwtConfig,
	)

	// Инициализируем конфигурацию
	if err := cfg.InitConfig(); err != nil {
		return nil, fmt.Errorf("ошибка инициализации конфигурации: %w", err)
	}

	return cfg, nil
}

