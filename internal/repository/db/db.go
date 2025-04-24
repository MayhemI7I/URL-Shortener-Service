package db

import (
	"fmt"
	"sync"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/infrastructure"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/repository"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/repository/postgres"
)

// DB представляет синглтон для работы с базой данных
type DB struct {
	config infrastructure.DBConfigProvider
	jwtConfig infrastructure.JWTConfigProvider
	conn   *postgres.DB
	mu     sync.RWMutex
}

var (
	instance *DB
	once     sync.Once
)

// GetInstance возвращает единственный экземпляр DB (Синглтон)
func GetInstance() (*DB, error) {
	var err error
	once.Do(func() {
		instance = &DB{}
	})
	return instance, err
}

// Init инициализирует базу данных
func (d *DB) Init(config infrastructure.DBConfigProvider, jwtConfig infrastructure.JWTConfigProvider) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.config = config
	d.jwtConfig = jwtConfig

	// Инициализация соединения с базой данных
	conn, err := postgres.NewDB(config.GetDSN())
	if err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	d.conn = conn
	return nil
}

// GetURLRepository возвращает репозиторий для работы с URL
func (d *DB) GetURLRepository() repository.URLRepository {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return postgres.NewPostgresURLStorage(d.conn)
}

// GetUserRepository возвращает репозиторий для работы с пользователями
func (d *DB) GetUserRepository() repository.UserRepository {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return postgres.NewPostgresUserStorage(d.conn, d.jwtConfig)
}

// Close закрывает соединение с базой данных
func (d *DB) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
