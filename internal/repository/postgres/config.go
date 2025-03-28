package database

import (
	"fmt"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/config"
)

// DBConfig конфигурация для базы данных
type DBConfig struct {
	config.BaseConfig
	DSN             string
	PoolSize        int
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	Driver          string
}

// NewDBConfig создает новую конфигурацию базы данных
func NewDBConfig() *DBConfig {
	return &DBConfig{
		BaseConfig:      config.NewBaseConfig("DB"),
		DSN:             "postgres://postgres:1@localhost:5432/usvideos",
		PoolSize:        10,
		MaxIdleConns:    5,
		MaxOpenConns:    25,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
		Driver:          "pgx",
	}
}

// AddFlags добавляет флаги для конфигурации базы данных
func (c *DBConfig) AddFlags() {
	// Используем флаги из BaseConfig
	c.AddStringFlag(c.GetFlags(), "dsn", c.DSN, "PostgreSQL DSN")
	c.AddIntFlag(c.GetFlags(), "pool-size", c.PoolSize, "Database connection pool size")
	c.AddIntFlag(c.GetFlags(), "max-idle-conns", c.MaxIdleConns, "Maximum number of idle database connections")
	c.AddIntFlag(c.GetFlags(), "max-open-conns", c.MaxOpenConns, "Maximum number of open database connections")
	c.AddDurationFlag(c.GetFlags(), "conn-max-lifetime", c.ConnMaxLifetime, "Maximum lifetime of database connections")
	c.AddDurationFlag(c.GetFlags(), "conn-max-idle-time", c.ConnMaxIdleTime, "Maximum idle time of database connections")
	c.AddStringFlag(c.GetFlags(), "driver", c.Driver, "Database driver")
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func (c *DBConfig) LoadFromEnv() {
	// Базовая реализация уже обрабатывает переменные окружения
	// через методы Add*Flag
}

// Validate проверяет корректность конфигурации
func (c *DBConfig) Validate() error {
	// Проверяем обязательные поля
	if err := c.ValidateRequired(map[string]interface{}{
		"dsn":    c.DSN,
		"driver": c.Driver,
	}); err != nil {
		return err
	}

	// Проверяем числовые значения
	if err := c.ValidatePositive(map[string]int{
		"pool-size":      c.PoolSize,
		"max-idle-conns": c.MaxIdleConns,
		"max-open-conns": c.MaxOpenConns,
	}); err != nil {
		return err
	}

	// Проверяем длительности
	if err := c.ValidateDuration(map[string]time.Duration{
		"conn-max-lifetime":  c.ConnMaxLifetime,
		"conn-max-idle-time": c.ConnMaxIdleTime,
	}); err != nil {
		return err
	}

	// Проверяем логические ограничения
	if c.MaxOpenConns < c.MaxIdleConns {
		return fmt.Errorf("максимальное количество открытых соединений должно быть больше или равно количеству неактивных соединений")
	}

	return nil
}

// GetDefaultConfig возвращает конфигурацию по умолчанию
func (c *DBConfig) GetDefaultConfig() interface{} {
	return NewDBConfig()
}
