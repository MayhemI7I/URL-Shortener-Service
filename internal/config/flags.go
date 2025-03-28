package config

import (
	"os"
	"strconv"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/auth"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/database"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"
	"github.com/spf13/pflag"
)

type FileStorage struct {
	Path string
}

// Config представляет конфигурацию приложения
type Config struct {
	ServerAdress    string
	ServerPort      string
	BaseURL         string
	Logger          logger.LoggerConfig
	FileStorage     FileStorage
	DB              database.DBConfig
	URLLength       uint16
	JWT             auth.JWTConfig
	ShutdownTimeout string
}

// InitConfig инициализирует конфигурацию приложения
func InitConfig() *Config {
	cfg := &Config{
		Logger: logger.DefaultLoggerConfig(),
		DB:     database.DefaultDBConfig(),
		JWT:    auth.DefaultJWTConfig(),
	}

	// Добавляем флаги для всех компонентов
	cfg.Logger.AddFlags()
	cfg.DB.AddFlags()
	cfg.JWT.AddFlags()

	// Загружаем конфигурацию из переменных окружения
	cfg.Logger.LoadFromEnv()
	cfg.DB.LoadFromEnv()
	cfg.JWT.LoadFromEnv()

	// Проверяем корректность конфигурации
	if err := cfg.Logger.Validate(); err != nil {
		logger.Log.Fatalf("Ошибка валидации конфигурации логгера: %v", err)
	}
	if err := cfg.DB.Validate(); err != nil {
		logger.Log.Fatalf("Ошибка валидации конфигурации базы данных: %v", err)
	}
	if err := cfg.JWT.Validate(); err != nil {
		logger.Log.Fatalf("Ошибка валидации конфигурации JWT: %v", err)
	}

	// Определяем флаги командной строки
	pflag.StringVarP(&cfg.ServerAdress, "server-address", "s", "localhost", "Server address")
	pflag.StringVarP(&cfg.ServerPort, "server-port", "p", "8080", "Server port")
	pflag.StringVarP(&cfg.BaseURL, "base-url", "b", "http://localhost:8080", "Base URL for return server")
	pflag.StringVarP(&cfg.FileStorage.Path, "file-storage", "f", "short-url-db.json", "Path to file storage")
	pflag.Uint16VarP(&cfg.URLLength, "url-length", "l", 8, "URL length")
	pflag.StringVarP(&cfg.ShutdownTimeout, "shutdown-timeout", "t", "10", "Shutdown timeout")

	// Переопределяем конфигурацию переменными окружения
	if envServerAdress := os.Getenv("SERVER_ADDRESS"); envServerAdress != "" {
		cfg.ServerAdress = envServerAdress
	}
	if envServerPort := os.Getenv("SERVER_PORT"); envServerPort != "" {
		cfg.ServerPort = envServerPort
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseURL = envBaseURL
	}
	if envFileStorage := os.Getenv("FILE_STORAGE"); envFileStorage != "" {
		cfg.FileStorage.Path = envFileStorage
	}
	if envURLLength := os.Getenv("URL_LENGTH"); envURLLength != "" {
		if urlLength, err := strconv.ParseUint(envURLLength, 10, 16); err == nil {
			cfg.URLLength = uint16(urlLength)
		}
	}
	if envShutdownTimeout := os.Getenv("SHUTDOWN_TIMEOUT"); envShutdownTimeout != "" {
		cfg.ShutdownTimeout = envShutdownTimeout
	}

	// Парсим флаги командной строки
	pflag.Parse()

	return cfg
}
