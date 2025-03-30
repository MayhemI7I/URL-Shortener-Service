package repository

import (
	"errors"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/config"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/repository/file"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/repository/memory"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/repository/postgres"

)

// StorageType определяет тип хранилища
type StorageType string

const (
	StorageTypePostgres StorageType = "postgres"
	StorageTypeMemory   StorageType = "memory"
	StorageTypeFile     StorageType = "file"
)

// StorageSelector выбирает подходящее хранилище на основе конфигурации
type StorageSelector struct {
	cfg *config.AppConfig
	db  *postgres.DB
}

// NewStorageSelector создает новый селектор хранилища
func NewStorageSelector(cfg *config.AppConfig) *StorageSelector {
	return &StorageSelector{cfg: cfg}
}

// SelectURLStorage выбирает и инициализирует хранилище URL наосновании переданной конфигурации при инициализации приложения
func (s *StorageSelector) SelectURLStorage() (interfaces.URLRepository, error) {
	switch s.getStorageType() {
	case StorageTypePostgres:
		if s.db == nil {
			db, err := postgres.NewDB(s.cfg.DBConfig.GetDSN())
			if err != nil {
				return nil, err
			}
			s.db = db
		}
		return postgres.NewPostgresURLStorage(s.db), nil
	case StorageTypeFile:
		return file.NewFileStorage(s.cfg.FileStorageConfig.GetPath())
	case StorageTypeMemory:
		return memory.NewMemoryStorage()
	default:
		return memory.NewMemoryStorage() // По умолчанию используем in-memory хранилище
	}
}

// SelectUserStorage выбирает и инициализирует хранилище пользователей
func (s *StorageSelector) SelectUserStorage() (interfaces.UserRepository, error) {
	switch s.getStorageType() {
	case StorageTypePostgres:
		if s.db == nil {
			db, err := postgres.NewDB(s.cfg.DBConfig.GetDSN())
			if err != nil {
				return nil, err
			}
			s.db = db
		}
		return postgres.NewPostgresUserStorage(s.db, s.cfg.JWTConfig), nil
	default:
		return nil, errors.New("user storage is only supported with PostgreSQL")
	}
}

// Close закрывает соединение с базой данных
func (s *StorageSelector) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// getStorageType определяет тип хранилища на основе конфигурации
func (s *StorageSelector) getStorageType() StorageType {
	if s.cfg.DBConfig.GetDriver() != "" {
		return StorageTypePostgres
	}
	if s.cfg.FileStorageConfig.GetPath() != "" {
		return StorageTypeFile
	}
	return StorageTypeMemory
}
