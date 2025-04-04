package infrastructure

import (
	"time"
)

// Configurable интерфейс для всех конфигураций
type Configurable interface {
	AddFlags()
	LoadFromEnv()
	Validate() error
}

// ConfigContainer интерфейс для контейнера конфигураций
type ConfigContainer interface {
	AddFlags()
	LoadFromEnv()
	Validate() error
	ParseFlags() error
	PrintHelp()
}

// HTTPConfigProvider интерфейс для HTTP конфигурации
type HTTPConfigProvider interface {
	Configurable
	GetAddress() string
	GetPort() string
	GetProtocol() string
	GetReadTimeout() int
	GetWriteTimeout() int
	GetShutdownTimeout() int
}

// DBConfigProvider интерфейс для конфигурации базы данных
type DBConfigProvider interface {
	Configurable
	GetDSN() string
	GetDriver() string
	GetMaxOpenConns() int
	GetMaxIdleConns() int
	GetConnMaxLifetime() int
}

type FileStorageConfigProvider interface {
	Configurable
	GetPath() string
}

// LoggerConfigProvider определяет интерфейс для конфигурации логгера
type LoggerConfigProvider interface {
	Configurable
	// GetLevel возвращает уровень логирования
	GetLevel() string

	// GetOutput возвращает тип вывода логов (file, console, both)
	GetOutput() string

	// GetLogPath возвращает путь к файлу логов
	GetLogPath() string

	// GetMaxSize возвращает максимальный размер файла логов в МБ
	GetMaxSize() int

	// GetMaxBackups возвращает максимальное количество резервных копий
	GetMaxBackups() int

	// GetMaxAge возвращает максимальный возраст файла логов в днях
	GetMaxAge() int

	// GetCompress возвращает флаг сжатия старых логов
	GetCompress() bool
}

// JWTConfigProvider интерфейс для JWT конфигурации
type JWTConfigProvider interface {
	Configurable
	GetSecretKey() string                     // секретный ключ для JWT токенов
	GetAccessTokenExpiration() time.Duration  // время жизни access-токена в секундах
	GetRefreshTokenExpiration() time.Duration // время жизни refresh-токена в секундах
}
