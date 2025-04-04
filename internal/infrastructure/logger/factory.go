package logger

import (
	"fmt"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/infrastructure"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger/adapters"
)

// LoggerType определяет тип логгера
type LoggerType string

const (
	ZapLoggerType LoggerType = "zap"
	// Можно добавить другие типы логгеров
)

// LoggerFactory создает логгеры разных типов
type LoggerFactory struct {
	config infrastructure.LoggerConfigProvider
}

// NewLoggerFactory создает новую фабрику логгеров
func NewLoggerFactory(config infrastructure.LoggerConfigProvider) *LoggerFactory {
	return &LoggerFactory{
		config: config,
	}
}

// CreateLogger создает логгер указанного типа
func (f *LoggerFactory) CreateLogger(loggerType LoggerType) (infrastructure.Logger, error) {
	switch loggerType {
	case ZapLoggerType:
		return adapters.NewZapLogger(f.config)
	default:
		return nil, fmt.Errorf("неподдерживаемый тип логгера: %s", loggerType)
	}
}

// CreateDefaultLogger создает логгер по умолчанию
func (f *LoggerFactory) CreateDefaultLogger() (infrastructure.Logger, error) {
	return f.CreateLogger(ZapLoggerType)
}
