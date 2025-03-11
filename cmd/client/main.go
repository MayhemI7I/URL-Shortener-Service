package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
	"net/http"

	"github.com/MayhemI7I/URL-Shortener-Service/config"
	"github.com/MayhemI7I/URL-Shortener-Service/domain"
	"github.com/MayhemI7I/URL-Shortener-Service/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// Константы для путей API
const (
	APIPathPing    = "/Ph-VaNhL"      // Проверка доступности сервера
	APIPathLogin   = "/api/login"     // Аутентификация
	APIPathShorten = "/api/shorten"   // Создание коротких URL через JSON
	APIPathRoot    = "/"              // Создание коротких URL через FormData
	APIGetAllURLs    = "/api/user/urls" // Получение всех URL пользователя
)

// ClientReq — клиент для HTTP-запросов с поддержкой кук
type ClientReq struct {
	request *resty.Client
}

// NewClientReq создаёт новый экземпляр клиента с таймаутом и CookieJar
func NewClientReq() (*ClientReq, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := resty.New().
		SetTimeout(500 * time.Millisecond).
		SetCookieJar(jar)

	return &ClientReq{request: client}, nil
}

// Login выполняет вход и получает куки
func (c *ClientReq) Login(url string) error {
	resp, err := c.request.R().Post(url)
	if err != nil {
		logger.Log.Error("failed to login", zap.Error(err))
		return fmt.Errorf("login request failed: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Log.Error("login failed with status", zap.Int("status", resp.StatusCode()))
		return fmt.Errorf("login failed with status: %d", resp.StatusCode())
	}

	// Выводим куки для отладки
	for _, cookie := range resp.Cookies() {
		fmt.Printf("Cookie: %s = %s\n", cookie.Name, cookie.Value)
	}
	logger.Log.Info("login successful")
	return nil
}

// GetAllURLs получает все URL пользователя
func (c *ClientReq) GetAllURLs(url string) ([]domain.URLData, error) {
	var datas []domain.URLData
	resp, err := c.request.R().Get(url)
	if err != nil {
		logger.Log.Error("failed to get all URLs", zap.Error(err))
		return nil, fmt.Errorf("get all URLs failed: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Log.Error("get all URLs failed with status", zap.Int("status", resp.StatusCode()))
		return nil, fmt.Errorf("get all URLs failed with status: %d", resp.StatusCode())
	}

	if err := json.Unmarshal(resp.Body(), &datas); err != nil {
		logger.Log.Error("failed to unmarshal URLs", zap.Error(err))
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	for _, data := range datas {
		fmt.Printf("User ID: %s\nShort URL: %s\nOriginal URL: %s\nCreated At: %s\n\n",
			data.ID, data.ShortURL, data.OrigURL, data.CreatedAt.Format(time.RFC3339))
	}
	return datas, nil
}

// PostJSON отправляет JSON-запрос
func (c *ClientReq) PostJSON(url string, body []byte) (string, int, error) {
	resp, err := c.request.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
	if err != nil {
		logger.Log.Error("post JSON failed", zap.Error(err))
		return "", 0, fmt.Errorf("post JSON request failed: %w", err)
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		logger.Log.Error("post JSON returned error", zap.Int("status", resp.StatusCode()), zap.String("body", resp.String()))
		return "", resp.StatusCode(), fmt.Errorf("server error: %s (status: %d)", resp.String(), resp.StatusCode())
	}
	return resp.String(), resp.StatusCode(), nil
}

// PostFormData отправляет FormData-запрос
func (c *ClientReq) PostFormData(url, origURL string) (string, int, error) {
	resp, err := c.request.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{"url": origURL}).
		Post(url)
	if err != nil {
		logger.Log.Error("post FormData failed", zap.Error(err))
		return "", 0, fmt.Errorf("post FormData request failed: %w", err)
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		logger.Log.Error("post FormData returned error", zap.Int("status", resp.StatusCode()), zap.String("body", resp.String()))
		return "", resp.StatusCode(), fmt.Errorf("server error: %s (status: %d)", resp.String(), resp.StatusCode())
	}
	return resp.String(), resp.StatusCode(), nil
}

// GetPing проверяет доступность сервера
func (c *ClientReq) GetPing(url string) (string, int, error) {
	resp, err := c.request.R().Get(url)
	if err != nil {
		logger.Log.Error("ping failed", zap.Error(err))
		return "", 0, fmt.Errorf("ping request failed: %w", err)
	}
	return resp.String(), resp.StatusCode(), nil
}

// GetShortURL проверяет редирект по короткому URL
func (c *ClientReq) GetShortURL(url string) (string, int, error) {
	resp, err := c.request.R().Get(url)
	if err != nil {
		logger.Log.Error("get short URL failed", zap.Error(err))
		return "", 0, fmt.Errorf("get short URL request failed: %w", err)
	}
	return resp.String(), resp.StatusCode(), nil
}

// NewUrlData создаёт новый экземпляр URLData
func NewUrlData(origURL string) *domain.URLData {
	return &domain.URLData{
		URLPair: domain.URLPair{
			OrigURL:   origURL,
			CreatedAt: time.Now(),
		},
	}
}

// runTests запускает все тесты клиента
func runTests(client *ClientReq, cfg *config.Config, urls []domain.URLData) {
	fmt.Println("\n=== ТЕСТИРУЕМ СЕРВЕР ===")

	fmt.Println("\n🔹 Тест: POST /api/login")
	if err := client.Login(cfg.BaseURL + APIPathLogin); err != nil {
		fmt.Printf("❌ Ошибка логина: %v\n", err)
	} else {
		fmt.Println("✅ Успешный логин")
	}

	
	fmt.Println("\n🔹 Тест: GET /ping")
	resp, status, err := client.GetPing(cfg.BaseURL + APIPathPing)
	if err == nil {
		fmt.Printf("✅ Сервер доступен! Ответ: %s (Код: %d)\n", resp, status)
	} else {
		fmt.Printf("❌ Ошибка ping: %v\n", err)
	}



	// 3. Тест POST /api/shorten (JSON)
	fmt.Println("\n🔹 Тест: POST /api/shorten (JSON)")
	data, err := json.Marshal(urls)
	if err != nil {
		fmt.Printf("❌ Ошибка маршалинга JSON: %v\n", err)
	} else {
		resp, status, err := client.PostJSON(cfg.BaseURL+APIPathShorten, data)
		if err == nil {
			fmt.Printf("✅ Короткие URL созданы: %s (Код: %d)\n", resp, status)
		} else {
			fmt.Printf("❌ Ошибка создания JSON URL: %v\n", err)
		}
	}

	// 4. Тест POST / (FormData)
	fmt.Println("\n🔹 Тест: POST / (FormData)")
	for _, url := range urls {
		resp, status, err := client.PostFormData(cfg.BaseURL+APIPathRoot, url.OrigURL)
		if err == nil {
			fmt.Printf("✅ Короткий URL создан: %s (Код: %d)\n", resp, status)
		} else {
			fmt.Printf("❌ Ошибка создания FormData URL: %v\n", err)
		}
	}

	//5 Тест GET всех URL 
	fmt.Println("\n🔹 Тест GET всех URL")

responseAllURLs, err := client.GetAllURLs(cfg.BaseURL + APIGetAllURLs)
if err != nil {
	fmt.Printf("❌ Ошибка получения всех URL: %v\n", err)
	return
}


countRecords := len(responseAllURLs)
fmt.Printf("✅ Все URL получены. Количество записей: %d\n", countRecords)

for _, url := range responseAllURLs {
	fmt.Printf("Для user: `%s`\nКоличество записей: %d \nОригинальный URL: %s\nСокращенный URL: %s\n\n", url.ID, countRecords, url.OrigURL, url.ShortURL)
}

	// 5. Тест GET всех URL (опционально, если есть такой endpoint)
	// fmt.Println("\n🔹 Тест: GET /api/user/urls")
	// if urls, err := client.GetAllURLs(cfg.BaseURL + "/api/user/urls"); err == nil {
	// 	fmt.Printf("✅ Получено %d URL\n", len(urls))
	// } else {
	// 	fmt.Printf("❌ Ошибка получения всех URL: %v\n", err)
	// }

	fmt.Println("\n✅ Все тесты завершены!")
}

func main() {
	// Инициализация конфигурации и логгера
	cfg := config.InitConfig()
	logger.InitLogger(cfg.LogLevel)
	defer logger.CloseLogger()

	logger.Log.Info("starting client")

	// Создание клиента
	client, err := NewClientReq()
	if err != nil {
		logger.Log.Fatal("failed to create client", zap.Error(err))
	}

	// Чтение количества URL
	fmt.Print("Введите количество URL для отправки: ")
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		logger.Log.Fatal("failed to read number of URLs", zap.Error(err))
	}
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n') // Пропускаем оставшийся \n

	// Чтение длинных URL
	urls := make([]domain.URLData, 0, n)
	for i := 0; i < n; i++ {
		fmt.Printf("Введите длинный URL #%d: ", i+1)
		origURL, err := reader.ReadString('\n')
		if err != nil {
			logger.Log.Fatal("failed to read URL", zap.Error(err))
		}
		origURL = strings.TrimSpace(origURL)
		urls = append(urls, *NewUrlData(origURL))
	}

	// Запуск тестов
	runTests(client, cfg, urls)
}