package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/auth"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/database"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"
	"github.com/spf13/pflag"
)

// AppConfig основная конфигурация приложения
type AppConfig struct {
	// Основные настройки приложения
	ServerAddr      string
	ServerPort      int
	BaseURL         string
	ShutdownTimeout time.Duration

	// Конфигурации компонентов
	Logger *logger.LoggerConfig
	DB     *database.DBConfig
	JWT    *auth.JWTConfig
	HTTP   *http.HTTPConfig

	// Флаги командной строки
	flags *pflag.FlagSet
}

// NewAppConfig создает новую конфигурацию приложения
func NewAppConfig() *AppConfig {
	return &AppConfig{
		ServerAddr:      "localhost",
		ServerPort:      8080,
		BaseURL:         "http://localhost:8080",
		ShutdownTimeout: 10 * time.Second,
		Logger:          logger.NewLoggerConfig(),
		DB:              database.NewDBConfig(),
		JWT:             auth.NewJWTConfig(),
		flags:           pflag.NewFlagSet("app", pflag.ExitOnError),
	}
}

// AddFlags добавляет все флаги командной строки
func (c *AppConfig) AddFlags() {
	// Основные флаги приложения
	c.flags.StringVarP(&c.ServerAddr, "addr", "a", c.ServerAddr, "Адрес сервера")
	c.flags.IntVarP(&c.ServerPort, "port", "p", c.ServerPort, "Порт сервера")
	c.flags.StringVarP(&c.BaseURL, "base-url", "b", c.BaseURL, "Базовый URL сервиса")
	c.flags.DurationVarP(&c.ShutdownTimeout, "shutdown-timeout", "t", c.ShutdownTimeout, "Таймаут для graceful shutdown")

	// Флаги компонентов
	c.Logger.AddFlags()
	c.DB.AddFlags()
	c.JWT.AddFlags()
	c.HTTP.AddFlags()
}

// ParseFlags парсит флаги командной строки
func (c *AppConfig) ParseFlags() error {
	return c.flags.Parse(os.Args[1:])
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func (c *AppConfig) LoadFromEnv() {
	// Загрузка основных настроек
	if addr := os.Getenv("SERVER_ADDR"); addr != "" {
		c.ServerAddr = addr
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.ServerPort = p
		}
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		c.BaseURL = baseURL
	}
	if timeout := os.Getenv("SHUTDOWN_TIMEOUT"); timeout != "" {
		if t, err := time.ParseDuration(timeout); err == nil {
			c.ShutdownTimeout = t
		}
	}

	// Загрузка конфигураций компонентов
	c.Logger.LoadFromEnv()
	c.DB.LoadFromEnv()
	c.JWT.LoadFromEnv()
	c.HTTP.LoadFromEnv()
}

// Validate проверяет корректность всей конфигурации
func (c *AppConfig) Validate() error {
	// Проверка основных настроек
	if c.ServerAddr == "" {
		return fmt.Errorf("адрес сервера не может быть пустым")
	}
	if c.ServerPort <= 0 {
		return fmt.Errorf("порт сервера должен быть положительным числом")
	}
	if c.BaseURL == "" {
		return fmt.Errorf("базовый URL не может быть пустым")
	}
	if c.ShutdownTimeout <= 0 {
		return fmt.Errorf("таймаут shutdown должен быть положительной длительностью")
	}

	// Проверка конфигураций компонентов
	if err := c.Logger.Validate(); err != nil {
		return fmt.Errorf("ошибка конфигурации логгера: %w", err)
	}
	if err := c.DB.Validate(); err != nil {
		return fmt.Errorf("ошибка конфигурации базы данных: %w", err)
	}
	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("ошибка конфигурации JWT: %w", err)
	}
	if err := c.HTTP.Validate(); err != nil {
		return fmt.Errorf("ошибка конфигурации HTTP: %w", err)
	}

	return nil
}

// GetServerAddr возвращает полный адрес сервера
func (c *AppConfig) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.ServerAddr, c.ServerPort)
}

// PrintHelp выводит справку по использованию флагов
func (c *AppConfig) PrintHelp() {
	c.flags.PrintDefaults()
}
