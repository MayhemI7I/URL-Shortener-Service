package logger

import (
	"fmt"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/config"
	"github.com/spf13/pflag"
)

// LoggerConfig конфигурация для логгера
type LoggerConfig struct {
	name        string
	flags       *pflag.FlagSet
	LogPath     string
	MaxSize     int
	MaxBackups  int
	MaxAge      int
	Compress    bool
	Level       string
	OutputType  string // "file", "console", "both"
}

// NewLoggerConfig создает новую конфигурацию логгера
func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		name:        "LOG",
		flags:       pflag.NewFlagSet("LOG", pflag.ExitOnError),
		LogPath:     "logs/app.log",
		MaxSize:     100,    // 100 МБ
		MaxBackups:  3,      // 3 резервные копии
		MaxAge:      30,     // 30 дней
		Compress:    true,   // сжатие старых логов
		Level:       "info", // уровень логирования по умолчанию
		OutputType:  "file", // тип вывода по умолчанию
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
	c.flags.StringVarP(&c.OutputType, "log-output", "", c.OutputType, "Тип вывода логов (file, console, both)")
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func (c *LoggerConfig) LoadFromEnv() {
	// Здесь можно добавить логику загрузки из env переменных
}

// Validate проверяет корректность конфигурации
func (c *LoggerConfig) Validate() error {
	// Проверяем обязательные поля
	if c.LogPath == "" && (c.OutputType == "file" || c.OutputType == "both") {
		return fmt.Errorf("поле log-path не может быть пустым при выводе в файл")
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

	// Проверяем тип вывода
	if !isValidOutputType(c.OutputType) {
		return fmt.Errorf("некорректный тип вывода логов: %s", c.OutputType)
	}

	return nil
}

// Получение значений

// GetLevel возвращает уровень логирования
func (c *LoggerConfig) GetLevel() string {
	return c.Level
}

// GetOutput возвращает тип вывода логов
func (c *LoggerConfig) GetOutput() string {
	return c.OutputType
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

// isValidOutputType проверяет корректность типа вывода
func isValidOutputType(output string) bool {
	validOutputs := map[string]bool{
		"file":    true,
		"console": true,
		"both":    true,
	}
	return validOutputs[output]
}

// Проверка соответствия интерфейсу
var _ config.LoggerConfigProvider = (*LoggerConfig)(nil)
var _ config.Configurable = (*LoggerConfig)(nil)
