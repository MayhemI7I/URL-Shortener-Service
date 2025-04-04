package http

import (
	"fmt"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/infrastructure"
	"github.com/spf13/pflag"
)

// HTTPConfig конфигурация HTTP сервера
type HTTPConfig struct {
	name            string // имя конфигурации
	flags           *pflag.FlagSet
	Address         string
	Port            string
	Protocol        string // http or https
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

// NewHTTPConfig создает новую конфигурацию HTTP
func NewHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		name:            "HTTP",
		flags:           pflag.NewFlagSet("HTTP", pflag.ExitOnError),
		Address:         "localhost",
		Port:            "8080",
		Protocol:        "http",
		ShutdownTimeout: 10 * time.Second,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
	}
}

// AddFlags добавляет флаги командной строки
func (c *HTTPConfig) AddFlags() {
	c.flags.StringVarP(&c.Address, "address", "", c.Address, "Address to listen on")
	c.flags.StringVarP(&c.Port, "port", "", c.Port, "Port to listen on")
	c.flags.StringVarP(&c.Protocol, "protocol", "", c.Protocol, "Protocol to listen on")
	c.flags.DurationVarP(&c.ReadTimeout, "read-timeout", "", c.ReadTimeout, "Read timeout")
	c.flags.DurationVarP(&c.WriteTimeout, "write-timeout", "", c.WriteTimeout, "Write timeout")
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func (c *HTTPConfig) LoadFromEnv() {
	// Можно добавить логику загрузки из env переменных
}

// Validate проверяет корректность конфигурации
func (c *HTTPConfig) Validate() error {
	// Проверяем обязательные поля
	if c.Address == "" {
		return fmt.Errorf("поле address не может быть пустым")
	}
	if c.Port == "" {
		return fmt.Errorf("поле port не может быть пустым")
	}
	if c.Protocol == "" {
		return fmt.Errorf("поле protocol не может быть пустым")
	}

	// Проверяем числовые значения
	if c.ReadTimeout <= 0 {
		return fmt.Errorf("поле read-timeout должно быть положительным числом")
	}
	if c.WriteTimeout <= 0 {
		return fmt.Errorf("поле write-timeout должно быть положительным числом")
	}

	return nil
}

// GetFlags возвращает набор флагов
func (c *HTTPConfig) GetFlags() *pflag.FlagSet {
	return c.flags
}

// Реализация HTTPConfigProvider

// GetAddress возвращает адрес сервера
func (c *HTTPConfig) GetAddress() string {
	return c.Address
}

// GetPort возвращает порт сервера
func (c *HTTPConfig) GetPort() string {
	return c.Port
}

// GetProtocol возвращает протокол сервера
func (c *HTTPConfig) GetProtocol() string {
	return c.Protocol
}

// GetReadTimeout возвращает таймаут чтения
func (c *HTTPConfig) GetReadTimeout() int {
	return int(c.ReadTimeout.Seconds())
}

// GetWriteTimeout возвращает таймаут записи
func (c *HTTPConfig) GetWriteTimeout() int {
	return int(c.WriteTimeout.Seconds())
}

// GetShutdownTimeout возвращает таймаут завершения работы
func (c *HTTPConfig) GetShutdownTimeout() int {
	return int(c.ShutdownTimeout.Seconds())
}

// Проверка соответствия интерфейсу
var _ infrastructure.HTTPConfigProvider = (*HTTPConfig)(nil)
var _ infrastructure.Configurable = (*HTTPConfig)(nil)





