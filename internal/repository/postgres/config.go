package postgres

import (
	"fmt"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/infrastructure"
	"github.com/spf13/pflag"
)

// DBConfig конфигурация для базы данных
type DBConfig struct {
	name            string
	flags           *pflag.FlagSet
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
		name:            "DB",
		flags:           pflag.NewFlagSet("DB", pflag.ExitOnError),
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
	c.flags.StringVarP(&c.DSN, "db-dsn", "", c.DSN, "PostgreSQL DSN")
	c.flags.IntVarP(&c.PoolSize, "db-pool-size", "", c.PoolSize, "Database connection pool size")
	c.flags.IntVarP(&c.MaxIdleConns, "db-max-idle-conns", "", c.MaxIdleConns, "Maximum number of idle database connections")
	c.flags.IntVarP(&c.MaxOpenConns, "db-max-open-conns", "", c.MaxOpenConns, "Maximum number of open database connections")
	c.flags.DurationVarP(&c.ConnMaxLifetime, "db-conn-max-lifetime", "", c.ConnMaxLifetime, "Maximum lifetime of database connections")
	c.flags.DurationVarP(&c.ConnMaxIdleTime, "db-conn-max-idle-time", "", c.ConnMaxIdleTime, "Maximum idle time of database connections")
	c.flags.StringVarP(&c.Driver, "db-driver", "", c.Driver, "Database driver")
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func (c *DBConfig) LoadFromEnv() {
	// Здесь можно добавить логику загрузки из env переменных
}

// Validate проверяет корректность конфигурации
func (c *DBConfig) Validate() error {
	// Проверяем обязательные поля
	if c.DSN == "" {
		return fmt.Errorf("поле db-dsn не может быть пустым")
	}
	if c.Driver == "" {
		return fmt.Errorf("поле db-driver не может быть пустым")
	}

	// Проверяем числовые значения
	if c.PoolSize <= 0 {
		return fmt.Errorf("поле db-pool-size должно быть положительным числом")
	}
	if c.MaxIdleConns <= 0 {
		return fmt.Errorf("поле db-max-idle-conns должно быть положительным числом")
	}
	if c.MaxOpenConns <= 0 {
		return fmt.Errorf("поле db-max-open-conns должно быть положительным числом")
	}

	// Проверяем длительности
	if c.ConnMaxLifetime <= 0 {
		return fmt.Errorf("поле db-conn-max-lifetime должно быть положительным числом")
	}
	if c.ConnMaxIdleTime <= 0 {
		return fmt.Errorf("поле db-conn-max-idle-time должно быть положительным числом")
	}

	// Проверяем логические ограничения
	if c.MaxOpenConns < c.MaxIdleConns {
		return fmt.Errorf("максимальное количество открытых соединений должно быть больше или равно количеству неактивных соединений")
	}

	return nil
}

// GetFlags возвращает набор флагов
func (c *DBConfig) GetFlags() *pflag.FlagSet {
	return c.flags
}

// GetDSN возвращает строку подключения к БД
func (c *DBConfig) GetDSN() string {
	return c.DSN
}

func (c *DBConfig) GetDriver() string {
	return c.Driver
}

// GetMaxOpenConns возвращает максимальное количество открытых соединений
func (c *DBConfig) GetMaxOpenConns() int {
	return c.MaxOpenConns
}

// GetMaxIdleConns возвращает максимальное количество неактивных соединений
func (c *DBConfig) GetMaxIdleConns() int {
	return c.MaxIdleConns
}

// GetConnMaxLifetime возвращает максимальное время жизни соединения в секундах
func (c *DBConfig) GetConnMaxLifetime() int {
	return int(c.ConnMaxLifetime.Seconds())
}

// Проверка соответствия интерфейсу
var _ infrastructure.DBConfigProvider = (*DBConfig)(nil)
var _ infrastructure.Configurable = (*DBConfig)(nil)
