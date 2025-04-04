package logger

import (
	"fmt"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/infrastructure"
	"github.com/spf13/pflag"
)

// LoggerConfig конфигурация для логгера
type LoggerConfig struct {
	name       string
	flags      *pflag.FlagSet
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
		name:       "LOG",
		flags:      pflag.NewFlagSet("LOG", pflag.ExitOnError),
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
	c.flags.StringVarP(&c.LogPath, "log-path", "", c.LogPath, "Путь к файлу логов")
	c.flags.IntVarP(&c.MaxSize, "log-max-size", "", c.MaxSize, "Максимальный размер файла логов в МБ")
	c.flags.IntVarP(&c.MaxBackups, "log-max-backups", "", c.MaxBackups, "Максимальное количество резервных копий")
	c.flags.IntVarP(&c.MaxAge, "log-max-age", "", c.MaxAge, "Максимальный возраст файла логов в днях")
	c.flags.BoolVarP(&c.Compress, "log-compress", "", c.Compress, "Сжимать старые логи")
	c.flags.StringVarP(&c.Level, "log-level", "", c.Level, "Уровень логирования (debug, info, warn, error)")
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func (c *LoggerConfig) LoadFromEnv() {
	// Здесь можно добавить логику загрузки из env переменных
}

// Validate проверяет корректность конфигурации
func (c *LoggerConfig) Validate() error {
	// Проверяем обязательные поля
	if c.LogPath == "" {
		return fmt.Errorf("поле log-path не может быть пустым")
	}
	if c.Level == "" {
		return fmt.Errorf("поле level не может быть пустым")
	}

	// Проверяем числовые значения
	if c.MaxSize <= 0 {
		return fmt.Errorf("поле max-size должно быть положительным числом")
	}
	if c.MaxBackups <= 0 {
		return fmt.Errorf("поле max-backups должно быть положительным числом")
	}
	if c.MaxAge <= 0 {
		return fmt.Errorf("поле max-age должно быть положительным числом")
	}

	// Проверяем уровень логирования
	if !isValidLogLevel(c.Level) {
		return fmt.Errorf("некорректный уровень логирования: %s", c.Level)
	}

	return nil
}

// GetLevel возвращает уровень логирования
func (c *LoggerConfig) GetLevel() string {
	return c.Level
}

// GetOutput всегда возвращает "both" для одновременного вывода в консоль и файл
func (c *LoggerConfig) GetOutput() string {
	return "both"
}

// GetLogPath возвращает путь к файлу логов
func (c *LoggerConfig) GetLogPath() string {
	return c.LogPath
}

// GetMaxSize возвращает максимальный размер файла логов в МБ
func (c *LoggerConfig) GetMaxSize() int {
	return c.MaxSize
}

// GetMaxBackups возвращает максимальное количество резервных копий
func (c *LoggerConfig) GetMaxBackups() int {
	return c.MaxBackups
}

// GetMaxAge возвращает максимальный возраст файла логов в днях
func (c *LoggerConfig) GetMaxAge() int {
	return c.MaxAge
}

// GetCompress возвращает флаг сжатия старых логов
func (c *LoggerConfig) GetCompress() bool {
	return c.Compress
}

// GetFlags возвращает набор флагов
func (c *LoggerConfig) GetFlags() *pflag.FlagSet {
	return c.flags
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

// Проверка соответствия интерфейсу
var _ infrastructure.LoggerConfigProvider = (*LoggerConfig)(nil)
var _ infrastructure.Configurable = (*LoggerConfig)(nil)
