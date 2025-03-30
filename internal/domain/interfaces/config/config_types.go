package config

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

// LoggerConfigProvider интерфейс для конфигурации логгера
type LoggerConfigProvider interface {
	Configurable
	GetLevel() string
	GetOutput() string
}

// JWTConfigProvider интерфейс для JWT конфигурации
type JWTConfigProvider interface {
	Configurable
	GetSecretKey() string // секретный ключ для JWT токенов
	GetTokenExpiration() int // время жизни токена в секундах
} 