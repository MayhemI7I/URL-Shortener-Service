package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/config"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/domain/interfaces/infrastructure"
	"github.com/MayhemI7I/URL-Shortener-Service/internal/usecases"
)

// App представляет основную структуру приложения
type App struct {
	config *config.AppConfig
	logger infrastructure.Logger
	core  *usecases.Core
}

// NewApp создает новый экземпляр приложения
func NewApp(cfg *config.AppConfig, log infrastructure.Logger, core *usecases.Core) *App {
	return &App{
		config: cfg,
		logger: log,
		core: core,
	}
}

// setupRouter настраивает маршрутизацию
func (a *App) setupRouter() http.Handler {
	// Создаем маршрутизатор
	mux := http.NewServeMux()

	// Настраиваем маршруты
	// Пример:
	// mux.Handle("/api/v1/endpoint", middleware.Chain(
	//    handler.NewEndpointHandler(a.service),
	//    middleware.Logger(a.logger),
	// ))

	// TODO: Добавьте здесь свои маршруты

	return mux
}

// Run запускает приложение
func (a *App) Run() error {
	// Настраиваем маршрутизатор
	router := a.setupRouter()

	// Создаем HTTP-сервер
	a.server = &http.Server{
		Addr:         a.config.GetServerAddr(),
		Handler:      router,
		ReadTimeout:  time.Duration(a.config.HTTPConfig.GetReadTimeout())
		WriteTimeout: time.Duration(a.)
		IdleTimeout:  60 * time.Second,
	}

	// Канал для ошибок сервера
	serverErrors := make(chan error, 1)

	// Запускаем сервер в отдельной горутине
	go func() {
		a.logger.Info("Сервер запущен",
			infrastructure.Field{Key: "address", Value: a.config.GetServerAddr()})
		serverErrors <- a.server.ListenAndServe()
	}()

	// Канал для сигналов ОС
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Блокируем до получения сигнала или ошибки
	select {
	case err := <-serverErrors:
		return fmt.Errorf("ошибка запуска сервера: %w", err)

	case sig := <-shutdown:
		a.logger.Info("Получен сигнал завершения работы",
			infrastructure.Field{Key: "signal", Value: sig.String()})

		// Создаем контекст с таймаутом для корректного завершения
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Пытаемся корректно завершить работу сервера
		if err := a.server.Shutdown(ctx); err != nil {
			// Если не получилось корректно завершить, делаем это принудительно
			a.server.Close()
			return fmt.Errorf("ошибка корректного завершения работы сервера: %w", err)
		}
	}

	return nil
}
