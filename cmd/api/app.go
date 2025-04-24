package api

import (
	"fmt"
	"sync"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/config"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/infrastructure"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/auth"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/repository/db"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/hash"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/http"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/logger"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/repository/file"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/repository/postgres"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/usecases"
)

// AppFacade представляет фасад для инициализации приложения
type AppFacade struct {
	config *config.AppConfig
	logger infrastructure.Logger
	core   *usecases.Core
	db     *db.DB
}

var (
	instance *AppFacade
	once     sync.Once
)

// GetInstance возвращает единственный экземпляр AppFacade (Синглтон)
func GetInstance() *AppFacade {
	once.Do(func() {
		instance = &AppFacade{}
	})
	return instance
}

// Init инициализирует приложение
func (a *AppFacade) Init() error {
	// Инициализация конфигурации
	if err := a.initConfig(); err != nil {
		return fmt.Errorf("ошибка инициализации конфигурации: %w", err)
	}

	// Инициализация логгера
	if err := a.initLogger(); err != nil {
		return fmt.Errorf("ошибка инициализации логгера: %w", err)
	}

	// Инициализация базы данных
	if err := a.initDB(); err != nil {
		return fmt.Errorf("ошибка инициализации базы данных: %w", err)
	}

	// Инициализация ядра приложения
	if err := a.initCore(); err != nil {
		return fmt.Errorf("ошибка инициализации ядра: %w", err)
	}

	return nil
}

// initConfig инициализирует конфигурацию
func (a *AppFacade) initConfig() error {
	// Создаем конфигурации компонентов
	httpConfig := http.NewHTTPConfig()
	loggerConfig := logger.NewLoggerConfig()
	dbConfig := postgres.NewDBConfig()
	fileStorageConfig := file.NewFileStorageConfig()
	jwtConfig := auth.NewJWTConfig()

	// Инициализируем конфигурацию приложения
	cfg, err := config.Init(
		httpConfig,
		loggerConfig,
		dbConfig,
		fileStorageConfig,
		jwtConfig,
	)
	if err != nil {
		return err
	}

	a.config = cfg
	return nil
}

// initLogger инициализирует логгер
func (a *AppFacade) initLogger() error {
	loggerFactory := logger.NewLoggerFactory(a.config.LoggerConfig)
	appLogger, err := loggerFactory.CreateDefaultLogger()
	if err != nil {
		return err
	}

	a.logger = appLogger
	return nil
}

// initDB инициализирует базу данных
func (a *AppFacade) initDB() error {
	dbInstance, err := db.GetInstance()
	if err != nil {
		return err
	}

	if err := dbInstance.Init(a.config.DBConfig, a.config.JWTConfig); err != nil {
		return err
	}

	a.db = dbInstance
	return nil
}

// initCore инициализирует ядро приложения
func (a *AppFacade) initCore() error {
	// Получаем репозитории из базы данных
	urlRepo := a.db.GetURLRepository()
	userRepo := a.db.GetUserRepository()

	// Создаем генератор URL
	urlGenerator := hash.NewHashGenerator()

	// Создаем ядро приложения
	core := usecases.NewCore(
		urlRepo,
		userRepo,
		urlGenerator,
		a.config.JWTConfig,
	)

	a.core = core
	return nil
}

// GetConfig возвращает конфигурацию приложения
func (a *AppFacade) GetConfig() *config.AppConfig {
	return a.config
}

// GetLogger возвращает логгер
func (a *AppFacade) GetLogger() infrastructure.Logger {
	return a.logger
}

// GetCore возвращает ядро приложения
func (a *AppFacade) GetCore() *usecases.Core {
	return a.core
}

// GetDB возвращает экземпляр базы данных
func (a *AppFacade) GetDB() *db.DB {
	return a.db
}
