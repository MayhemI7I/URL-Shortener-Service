package adapters

import (
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"

	"go.uber.org/zap"
)

// ZapLogger - реализация логгера на основе Zap
type ZapLogger struct {
	logger *zap.SugaredLogger
}

// NewZapLogger - создает новый логгер на основе Zap
func NewZapLogger(zapLogger *zap.SugaredLogger) interfaces.Logger {
	return &ZapLogger{
		logger: zapLogger,
	}
}

// Debug - реализует метод интерфейса Logger
func (z *ZapLogger) Debug(msg string, fields ...interfaces.Field) {
	args := convertFields(fields)
	z.logger.Debugw(msg, args...)
}

// Info - реализует метод интерфейса Logger
func (z *ZapLogger) Info(msg string, fields ...interfaces.Field) {
	args := convertFields(fields)
	z.logger.Infow(msg, args...)
}

// Warn - реализует метод интерфейса Logger
func (z *ZapLogger) Warn(msg string, fields ...interfaces.Field) {
	args := convertFields(fields)
	z.logger.Warnw(msg, args...)
}

// Error - реализует метод интерфейса Logger
func (z *ZapLogger) Error(msg string, fields ...interfaces.Field) {
	args := convertFields(fields)
	z.logger.Errorw(msg, args...)
}

// Fatal - реализует метод интерфейса Logger
func (z *ZapLogger) Fatal(msg string, fields ...interfaces.Field) {
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
