package domain

import (
	"time"

	"github.com/spf13/pflag"
)

// ConfigProvider определяет интерфейс для любого компонента,
// который может предоставить конфигурацию
type ConfigProvider interface {
	AddFlags()
	LoadFromEnv()
	Validate() error
	GetDefaultConfig() interface{}
}

// FlagProvider определяет интерфейс для работы с флагами
type FlagProvider interface {
	GetFlags() *pflag.FlagSet
	AddStringFlag(flags *pflag.FlagSet, name, value, usage string)
	AddIntFlag(flags *pflag.FlagSet, name string, value int, usage string)
	AddDurationFlag(flags *pflag.FlagSet, name string, value time.Duration, usage string)
	AddBoolFlag(flags *pflag.FlagSet, name string, value bool, usage string)
}

// ValidationProvider определяет интерфейс для валидации конфигурации
type ValidationProvider interface {
	ValidateRequired(fields map[string]interface{}) error
	ValidatePositive(fields map[string]int) error
	ValidateDuration(fields map[string]time.Duration) error
} 