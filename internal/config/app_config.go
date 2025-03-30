package config

import (
	"fmt"
	"os"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/config"
	"github.com/spf13/pflag"
)

// AppConfig основная конфигурация приложения
type AppConfig struct {
	// Конфигурации компонентов
	HTTPConfig config.HTTPConfigProvider
	LoggerConfig config.LoggerConfigProvider
	DBConfig config.DBConfigProvider
	FileStorageConfig config.FileStorageConfigProvider
	JWTConfig config.JWTConfigProvider


	// Флаги командной строки
	flags *pflag.FlagSet
}

// Создаем функцию-фабрику, которая примет конкретные реализации извне
func NewAppConfig(
	httpConfig config.HTTPConfigProvider,
	loggerConfig config.LoggerConfigProvider,
	dbConfig config.DBConfigProvider,
	fileStorageConfig config.FileStorageConfigProvider,
	jwtConfig config.JWTConfigProvider,
) *AppConfig {
	return &AppConfig{
		HTTPConfig: httpConfig,
		LoggerConfig: loggerConfig,
		DBConfig: dbConfig,
		FileStorageConfig: fileStorageConfig,
		JWTConfig: jwtConfig,
		flags: pflag.NewFlagSet("app", pflag.ExitOnError),
	}
}

// AddFlags добавляет все флаги командной строки
func (c *AppConfig) AddFlags() {
	// Флаги компонентов
	c.HTTPConfig.AddFlags()
	c.LoggerConfig.AddFlags()
	c.DBConfig.AddFlags()
	c.JWTConfig.AddFlags()
	c.FileStorageConfig.AddFlags()
}

// ParseFlags парсит флаги командной строки
func (c *AppConfig) ParseFlags() error {
	return c.flags.Parse(os.Args[1:])
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func (c *AppConfig) LoadFromEnv() {
	c.HTTPConfig.LoadFromEnv()
	c.LoggerConfig.LoadFromEnv()
	c.DBConfig.LoadFromEnv()
	c.JWTConfig.LoadFromEnv()
	c.FileStorageConfig.LoadFromEnv()
}

// Validate проверяет корректность всей конфигурации
func (c *AppConfig) Validate() error {
	// Проверка конфигураций компонентов
	if err := c.HTTPConfig.Validate(); err != nil {
		return fmt.Errorf("ошибка конфигурации HTTP: %w", err)
	}
	if err := c.LoggerConfig.Validate(); err != nil {
		return fmt.Errorf("ошибка конфигурации логгера: %w", err)
	}
	if err := c.DBConfig.Validate(); err != nil {
		return fmt.Errorf("ошибка конфигурации базы данных: %w", err)
	}
	if err := c.JWTConfig.Validate(); err != nil {
		return fmt.Errorf("ошибка конфигурации JWT: %w", err)
	}
	if err := c.FileStorageConfig.Validate(); err != nil {
		return fmt.Errorf("ошибка конфигурации файлового хранилища: %w", err)
	}
			

	return nil
}

// GetServerAddr возвращает полный адрес сервера
func (c *AppConfig) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", c.HTTPConfig.GetAddress(), c.HTTPConfig.GetPort())
}

// PrintHelp выводит справку по использованию флагов
func (c *AppConfig) PrintHelp() {
	c.flags.PrintDefaults()
}

// Проверка соответствия интерфейсу
var _ config.ConfigContainer = (*AppConfig)(nil)
