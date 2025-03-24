// Создаем конфигурацию для аутентификации
authConfig := config.NewAuthConfig(cfg)

// Создаем сервис аутентификации
authService := services.NewAuthService(userRepo, authConfig) 