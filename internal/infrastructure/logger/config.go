package logger

import (
	"fmt"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/config"
	"github.com/spf13/pflag"
)

// LoggerConfig конфигурация для логгера
type LoggerConfig struct {
	config.BaseConfig
	LogPath    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Level      string
}

// NewLoggerConfig создает новую конфигурацию логгера
func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		BaseConfig: config.NewBaseConfig("LOG"),
		LogPath:    "logs/app.log",
		MaxSize:    100,    // 100 МБ
		MaxBackups: 3,      // 3 резервные копии
		MaxAge:     30,     // 30 дней
		Compress:   true,   // сжатие старых логов
		Level:      "info", // уровень логирования по умолчанию
	}
}

// AddFlags добавляет флаги для конфигурации логгера
func (c *LoggerConfig) AddFlags() {
	flags := pflag.NewFlagSet("logger", pflag.ExitOnError)
	c.AddStringFlag(flags, "path", c.LogPath, "Путь к файлу логов")
	c.AddIntFlag(flags, "max-size", c.MaxSize, "Максимальный размер файла логов в МБ")
	c.AddIntFlag(flags, "max-backups", c.MaxBackups, "Максимальное количество резервных копий")
	c.AddIntFlag(flags, "max-age", c.MaxAge, "Максимальный возраст файла логов в днях")
	c.AddBoolFlag(flags, "compress", c.Compress, "Сжимать старые логи")
	c.AddStringFlag(flags, "level", c.Level, "Уровень логирования (debug, info, warn, error)")
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func (c *LoggerConfig) LoadFromEnv() {
	// Базовая реализация уже обрабатывает переменные окружения
	// через методы Add*Flag
}

// Validate проверяет корректность конфигурации
func (c *LoggerConfig) Validate() error {
	// Проверяем обязательные поля
	if err := c.ValidateRequired(map[string]interface{}{
		"path":  c.LogPath,
		"level": c.Level,
	}); err != nil {
		return err
	}

	// Проверяем числовые значения
	if err := c.ValidatePositive(map[string]int{
		"max-size":    c.MaxSize,
		"max-backups": c.MaxBackups,
		"max-age":     c.MaxAge,
	}); err != nil {
		return err
	}

	// Проверяем уровень логирования
	if !isValidLogLevel(c.Level) {
		return fmt.Errorf("некорректный уровень логирования: %s", c.Level)
	}

	return nil
}

// GetDefaultConfig возвращает конфигурацию по умолчанию
func (c *LoggerConfig) GetDefaultConfig() interface{} {
	return NewLoggerConfig()
}

// isValidLogLevel проверяет корректность уровня логирования
func isValidLogLevel(level string) bool {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	return validLevels[level]
}
