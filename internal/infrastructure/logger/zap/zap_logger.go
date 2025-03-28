package logger

import (
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger/adapters"
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

// InitLogger инициализирует базовый Zap логгер с ротацией логов
func InitLogger(config LoggerConfig) {
	adapter := adapters.NewZapAdapter(adapters.ZapLoggerConfig{
		LogPath:    config.LogPath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
		Level:      config.Level,
	})

	logger := adapter.AsLogger()
	Log = logger.(*adapters.ZapLoggerImpl).GetLogger()
}

// NewLogger создает новый логгер по умолчанию (ZAP)
func NewLogger() interfaces.Logger {
	if Log == nil {
		InitLogger(DefaultLoggerConfig())
	}
	return adapters.NewZapAdapter(adapters.ZapLoggerConfig{
		LogPath:    "logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
		Level:      "info",
	}).AsLogger()
}

// NewCustomLogger создает новый логгер указанного типа
func NewCustomLogger(loggerType string, customLogger interface{}) interfaces.Logger {
	switch loggerType {
	case "zap":
		if Log == nil {
			InitLogger(DefaultLoggerConfig())
		}
		return adapters.NewZapAdapter(adapters.ZapLoggerConfig{
			LogPath:    "logs/app.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     30,
			Compress:   true,
			Level:      "info",
		}).AsLogger()
	case "stdout":
		if adapter, ok := customLogger.(interfaces.LoggerAdapter); ok {
			return adapter.AsLogger()
		}
		return adapters.NewStdoutAdapter("info").AsLogger()
	default:
		// По умолчанию используем ZAP
		if Log == nil {
			InitLogger(DefaultLoggerConfig())
		}
		return adapters.NewZapAdapter(adapters.ZapLoggerConfig{
			LogPath:    "logs/app.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     30,
			Compress:   true,
			Level:      "info",
		}).AsLogger()
	}
}

// CloseLogger закрывает логгер
func CloseLogger() {
	if Log != nil {
		_ = Log.Sync()
	}
}
