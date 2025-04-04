package adapters

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/infrastructure"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapLogger реализация логгера на основе Zap
type ZapLogger struct {
	logger *zap.SugaredLogger
}

// NewZapLogger создает новый логгер на основе Zap
func NewZapLogger(config infrastructure.LoggerConfigProvider) (infrastructure.Logger, error) {
	// Устанавливаем уровень логирования
	level, err := getZapLevel(config.GetLevel())
	if err != nil {
		return nil, fmt.Errorf("ошибка установки уровня логирования: %w", err)
	}

	// Создаем конфигурацию энкодера
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Создаем энкодеры
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Создаем writers
	var cores []zapcore.Core

	// Всегда добавляем консольный вывод
	consoleWriter := zapcore.AddSync(zapcore.Lock(os.Stdout))
	cores = append(cores, zapcore.NewCore(consoleEncoder, consoleWriter, level))

	// Всегда добавляем файловый вывод
	if err := os.MkdirAll(filepath.Dir(config.GetLogPath()), 0755); err != nil {
		return nil, fmt.Errorf("ошибка создания директории для логов: %w", err)
	}

	// Настраиваем ротацию логов
	fileWriter := &lumberjack.Logger{
		Filename:   config.GetLogPath(),
		MaxSize:    config.GetMaxSize(),    // МБ
		MaxBackups: config.GetMaxBackups(), // количество файлов
		MaxAge:     config.GetMaxAge(),     // дней
		Compress:   config.GetCompress(),   // сжатие
	}

	cores = append(cores, zapcore.NewCore(jsonEncoder, zapcore.AddSync(fileWriter), level))

	// Создаем core
	core := zapcore.NewTee(cores...)

	// Создаем логгер
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()

	return &ZapLogger{
		logger: logger,
	}, nil
}

// Debug реализует метод интерфейса Logger
func (z *ZapLogger) Debug(msg string, fields ...infrastructure.Field) {
	args := convertFields(fields)
	z.logger.Debugw(msg, args...)
}

// Info реализует метод интерфейса Logger
func (z *ZapLogger) Info(msg string, fields ...infrastructure.Field) {
	args := convertFields(fields)
	z.logger.Infow(msg, args...)
}

// Warn реализует метод интерфейса Logger
func (z *ZapLogger) Warn(msg string, fields ...infrastructure.Field) {
	args := convertFields(fields)
	z.logger.Warnw(msg, args...)
}

// Error реализует метод интерфейса Logger
func (z *ZapLogger) Error(msg string, fields ...infrastructure.Field) {
	args := convertFields(fields)
	z.logger.Errorw(msg, args...)
}

// Fatal реализует метод интерфейса Logger
func (z *ZapLogger) Fatal(msg string, fields ...infrastructure.Field) {
	args := convertFields(fields)
	z.logger.Fatalw(msg, args...)
}



// convertFields конвертирует поля из типа interfaces.Field в аргументы для zap
func convertFields(fields []infrastructure.Field) []interface{} {
	if len(fields) == 0 {
		return nil
	}

	args := make([]interface{}, 0, len(fields)*2)
	for _, field := range fields {
		args = append(args, field.Key, field.Value)
	}

	return args
}

// getZapLevel конвертирует строковый уровень логирования в zapcore.Level
func getZapLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("неподдерживаемый уровень логирования: %s", level)
	}
}

// Проверка соответствия интерфейсу
var _ infrastructure.Logger = (*ZapLogger)(nil)

