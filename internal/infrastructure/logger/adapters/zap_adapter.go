package adapters

import (
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"

	"go.uber.org/zap"
)

// ZapLoggerAdapter - адаптер для Zap логгера
type ZapLoggerAdapter struct {
	logger *zap.SugaredLogger
}

// NewZapAdapter - создает новый адаптер для Zap логгера
func NewZapAdapter(zapLogger *zap.SugaredLogger) interfaces.LoggerAdapter {
	return &ZapLoggerAdapter{
		logger: zapLogger,
	}
}

// AsLogger - реализует интерфейс LoggerAdapter
func (z *ZapLoggerAdapter) AsLogger() interfaces.Logger {
	return &zapLoggerImpl{
		logger: z.logger,
	}
}

// zapLoggerImpl - реализация интерфейса Logger для Zap
type zapLoggerImpl struct {
	logger *zap.SugaredLogger
}

// Debug - реализует метод интерфейса Logger
func (z *zapLoggerImpl) Debug(msg string, fields ...interfaces.Field) {
	args := convertFields(fields)
	z.logger.Debugw(msg, args...)
}

// Info - реализует метод интерфейса Logger
func (z *zapLoggerImpl) Info(msg string, fields ...interfaces.Field) {
	args := convertFields(fields)
	z.logger.Infow(msg, args...)
}

// Warn - реализует метод интерфейса Logger
func (z *zapLoggerImpl) Warn(msg string, fields ...interfaces.Field) {
	args := convertFields(fields)
	z.logger.Warnw(msg, args...)
}

// Error - реализует метод интерфейса Logger
func (z *zapLoggerImpl) Error(msg string, fields ...interfaces.Field) {
	args := convertFields(fields)
	z.logger.Errorw(msg, args...)
}

// Fatal - реализует метод интерфейса Logger
func (z *zapLoggerImpl) Fatal(msg string, fields ...interfaces.Field) {
	args := convertFields(fields)
	z.logger.Fatalw(msg, args...)
}

// convertFields - конвертирует поля из типа interfaces.Field в аргументы для zap
func convertFields(fields []interfaces.Field) []interface{} {
	if len(fields) == 0 {
		return nil
	}

	args := make([]interface{}, 0, len(fields)*2)
	for _, field := range fields {
		args = append(args, field.Key, field.Value)
	}

	return args
}
