package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

// Configurable интерфейс для конфигураций 
type Configurable interface {
	AddFlags()
	LoadFromEnv()
	Validate() error
	GetDefaultConfig() interface{}
}

// BaseConfig базовый класс для конфигураций
type BaseConfig struct {
	name  string
	flags *pflag.FlagSet
}

// NewBaseConfig создает новую базовую конфигурацию
func NewBaseConfig(name string) BaseConfig {
	return BaseConfig{
		name:  name,
		flags: pflag.NewFlagSet(name, pflag.ExitOnError),
	}
}

// AddStringFlag добавляет строковый флаг
func (c *BaseConfig) AddStringFlag(flags *pflag.FlagSet, name, value, usage string) {
	flags.StringVarP(&value, name, "", value, usage)
}

// AddIntFlag добавляет целочисленный флаг
func (c *BaseConfig) AddIntFlag(flags *pflag.FlagSet, name string, value int, usage string) {
	flags.IntVarP(&value, name, "", value, usage)
}

// AddDurationFlag добавляет флаг длительности
func (c *BaseConfig) AddDurationFlag(flags *pflag.FlagSet, name string, value time.Duration, usage string) {
	flags.DurationVarP(&value, name, "", value, usage)
}

// AddBoolFlag добавляет булевый флаг
func (c *BaseConfig) AddBoolFlag(flags *pflag.FlagSet, name string, value bool, usage string) {
	flags.BoolVarP(&value, name, "", value, usage)
}

// ValidateRequired проверяет обязательные поля
func (c *BaseConfig) ValidateRequired(fields map[string]interface{}) error {
	for name, value := range fields {
		switch v := value.(type) {
		case string:
			if v == "" {
				return fmt.Errorf("поле %s не может быть пустым", name)
			}
		case int:
			if v == 0 {
				return fmt.Errorf("поле %s не может быть нулевым", name)
			}
		case time.Duration:
			if v == 0 {
				return fmt.Errorf("поле %s не может быть нулевым", name)
			}
		}
	}
	return nil
}

// ValidatePositive проверяет положительные значения
func (c *BaseConfig) ValidatePositive(fields map[string]int) error {
	for name, value := range fields {
		if value <= 0 {
			return fmt.Errorf("поле %s должно быть положительным числом", name)
		}
	}
	return nil
}

// ValidateDuration проверяет длительности
func (c *BaseConfig) ValidateDuration(fields map[string]time.Duration) error {
	for name, value := range fields {
		if value <= 0 {
			return fmt.Errorf("поле %s должно быть положительной длительностью", name)
		}
	}
	return nil
}

// LoadFromEnv загружает значение из переменной окружения
func (c *BaseConfig) LoadFromEnv(key string) string {
	return os.Getenv(key)
}

// LoadIntFromEnv загружает целое число из переменной окружения
func (c *BaseConfig) LoadIntFromEnv(key string) (int, error) {
	if value := os.Getenv(key); value != "" {
		return strconv.Atoi(value)
	}
	return 0, nil
}

// LoadDurationFromEnv загружает длительность из переменной окружения
func (c *BaseConfig) LoadDurationFromEnv(key string) (time.Duration, error) {
	if value := os.Getenv(key); value != "" {
		return time.ParseDuration(value)
	}
	return 0, nil
}

// GetFlags возвращает набор флагов
func (c *BaseConfig) GetFlags() *pflag.FlagSet {
	return c.flags
}

// PrintHelp выводит справку по использованию флагов
func (c *BaseConfig) PrintHelp() {
	c.flags.PrintDefaults()
}
