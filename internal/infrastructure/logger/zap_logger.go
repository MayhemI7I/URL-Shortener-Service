package logger

import (
	"os"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger/adapters"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

// InitLogger инициализирует базовый Zap логгер
func InitLogger(logLevel string) {
	// Устанавливаем уровень логирования
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	default:
		level = zapcore.InfoLevel
	}

	// Форматирование логов
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Для консоли
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleWriter := zapcore.AddSync(zapcore.Lock(os.Stdout))

	// Для файла
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Не удалось создать файл для логов: " + err.Error())
	}
	fileWriter := zapcore.AddSync(zapcore.Lock(logFile))

	// Создаем Core с консолью и файлом
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriter, level),
		zapcore.NewCore(fileEncoder, fileWriter, level),
	)

	// Создаем логгер
	logger := zap.New(core)
	Log = logger.Sugar()
}

// NewLogger создает новый логгер по умолчанию (ZAP)
func NewLogger() interfaces.Logger {
	if Log == nil {
		InitLogger("info")
	}
	return adapters.NewZapLogger(Log)
}

// NewCustomLogger создает новый логгер указанного типа
func NewCustomLogger(loggerType string, customLogger interface{}) interfaces.Logger {
	switch loggerType {
	case "zap":
		if Log == nil {
			InitLogger("info")
		}
		return adapters.NewZapLogger(Log)
	case "stdout":
		return adapters.NewStdoutLogger("info")
	default:
		// По умолчанию используем ZAP
		if Log == nil {
			InitLogger("info")
		}
		return adapters.NewZapLogger(Log)
	}
}

// CloseLogger закрывает логгер
func CloseLogger() {
	if Log != nil {
		_ = Log.Sync()
	}
}
