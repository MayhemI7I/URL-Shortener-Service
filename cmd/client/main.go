package main

import (
	"bufio"
	"encoding/json"
	"fmt"

	"local/config"
	"local/handlers/urlhandler"
	"local/logger"
	"os"
	"strings"
	"time"
	
	"github.com/go-resty/resty/v2"
)

// Клиент для HTTP-запросов
type ClientReq struct {
	request *resty.Client
}


// POST JSON (API)
func (c *ClientReq) PostJSON(url string, json []byte) (string, int, error) {
	response, err := c.request.R().
		SetHeader("Content-Type", "application/json").
		SetBody(json).
		Post(url)

	if err != nil {
		logger.Log.Errorf("error: %s", err.Error())
		return "", 0, err
	}

	if response.StatusCode() != 200 && response.StatusCode() != 201 {
		logger.Log.Errorf("Server returned error: %s, %v", response.Status(), response.String())
		return "", response.StatusCode(), fmt.Errorf("server error: %s", response.String())
	}

	return response.String(), response.StatusCode(), nil
}

// POST FormData (обычный POST-запрос)
func (c *ClientReq) PostFormData(url, longURL string) (string, int, error) {
	response, err := c.request.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{"url": longURL}).
		Post(url)

	if err != nil {
		logger.Log.Errorf("error: %s", err.Error())
		return "", 0, err
	}

	if response.StatusCode() != 200 && response.StatusCode() != 201 {
		logger.Log.Errorf("Server returned error: %s, %v", response.Status(), response.String())
		return "", response.StatusCode(), fmt.Errorf("server error: %s", response.String())
	}

	return response.String(), response.StatusCode(), nil
}

// GET /ping – проверка сервера
func (c *ClientReq) GetPing(url string) (string, int, error) {
	response, err := c.request.R().Get(url)

	if err != nil {
		logger.Log.Errorf("Ping error: %s", err.Error())
		return "", 0, err
	}

	return response.String(), response.StatusCode(), nil
}

// GET /{shortURL} – проверка редиректа
func (c *ClientReq) GetShortURL(url string) (string, int, error) {
	response, err := c.request.R().Get(url)

	if err != nil {
		logger.Log.Errorf("GET error: %s", err.Error())
		return "", 0, err
	}

	return response.String(), response.StatusCode(), nil
}

// Чтение длинного URL с консоли
func readLongURL() (string, error) {
	fmt.Println("Введите длинный URL:")
	reader := bufio.NewReader(os.Stdin)
	long, err := reader.ReadString('\n')
	long = strings.TrimSpace(long)
	if err != nil {
		return "", err
	}
	return long, nil
}

func main() {
	// Инициализация конфигурации и логгера
	cfg := config.InitConfig()
	logger.InitLogger(cfg.LogLevel)
	defer logger.CloseLogger()

	logger.Log.Info("Starting client")

	// Чтение длинного URL
	longURL, err := readLongURL()
	if err != nil {
		logger.Log.Fatalf("error: %s", err.Error())
	}
	r := urlhandler.NewURLRequest(longURL)


	// Создание HTTP клиента
	client := resty.New()
	client.SetTimeout(500 * time.Millisecond)
	postClient := &ClientReq{request: client}

	var response string
	var statusCode int

	// Тестируем все возможные пути
	fmt.Println("\n=== ТЕСТИРУЕМ СЕРВЕР ===")

	// 1. Проверяем, работает ли сервер (GET /ping)
	fmt.Println("\n🔹 Тест: GET /ping")
	response, statusCode, err = postClient.GetPing(cfg.BaseURL + "/Ph-VaNhL")
	if err == nil {
		fmt.Printf("✅ Сервер доступен! Ответ: %s (Код: %d)\n", response, statusCode)
	} else {
		fmt.Printf("❌ Ошибка ping: %v\n", err)
	}

	// 2. Тестируем создание короткого URL через JSON (POST /api/shorten)
	fmt.Println("\n🔹 Тест: POST /api/shorten (JSON)")

	n, err := json.Marshal(r)
	if err != nil {
		logger.Log.Fatalf("error: %s", err.Error())
	}
	response, statusCode, err = postClient.PostJSON(cfg.BaseURL+"/api/shorten",n )
	if err == nil {
		fmt.Printf("✅ Короткий URL создан: %s (Код: %d)\n", response, statusCode)
	} else {
		fmt.Printf("❌ Ошибка создания JSON URL: %v\n", err)
	}

	// 3. Тестируем создание короткого URL через FormData (POST /)
	fmt.Println("\n🔹 Тест: POST / (FormData)")
	response, statusCode, err = postClient.PostFormData(cfg.BaseURL+"/", longURL)
	if err == nil {
		fmt.Printf("✅ Короткий URL создан: %s (Код: %d)\n", response, statusCode)
	} else {
		fmt.Printf("❌ Ошибка создания FormData URL: %v\n", err)
	}

	// 4. Тестируем редирект по короткому URL (GET /{shortURL})
	// 4. Тестируем редирект по короткому URL (GET /{shortURL})
	if statusCode == 201 || statusCode == 200 {
		shortURL := strings.Trim(response, `"`) // Убираем кавычки, если сервер вернул JSON строку
		fmt.Println("\n🔹 Тест: GET redirect " + shortURL)

		resp, err := postClient.request.R().Get(cfg.BaseURL + "/" + shortURL)
		if err == nil && statusCode >= 300 && statusCode < 400 {
			location := resp.Header().Get("Location")
			if location != "" {
				fmt.Printf("✅ Редирект работает! Перенаправление на: %s (Код: %d)\n", location, statusCode)
			} else {
				fmt.Printf("❌ Ошибка: сервер не вернул заголовок Location (Код: %d)\n", statusCode)
			}
		} else {
			fmt.Printf("❌ Ошибка при редиректе: %v (Ответ: %s, Код: %d)\n", err, resp.String(), statusCode)
		}
	}

	// 5. Тестируем несуществующий URL (ошибочный запрос)
	fmt.Println("\n🔹 Тест: GET /unknown_path")
	response, statusCode, err = postClient.GetShortURL(cfg.BaseURL + "/unknown_path")
	if err == nil {
		fmt.Printf("✅ Сервер вернул ошибку ожидаемо: %s (Код: %d)\n", response, statusCode)
	} else {
		fmt.Printf("❌ Ошибка обработки неизвестного пути: %v\n", err)
	}

	fmt.Println("\n✅ Все тесты завершены!")
}
