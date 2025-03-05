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

	// Чтение количества объектов
	fmt.Print("Введите количество объектов для отправки: ")
	var n int
	fmt.Scan(&n)

	// Слайс для хранения всех запросов
	rs := make([]urlhandler.URLRequest, 0, n)
	reader := bufio.NewReader(os.Stdin)
    reader.ReadString('\n') // Читаем \n после ввода числа

	// Чтение длинных URL-ов
	for i := 0; i < n; i++ {
		longURL, err := reader.ReadString('\n')
		if err != nil {
			logger.Log.Fatalf("error: %s", err.Error())
		}
		r := urlhandler.NewURLRequest(longURL)
		rs = append(rs, *r) // Добавляем новый URLRequest в слайс
	}

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
	response, statusCode, err := postClient.GetPing(cfg.BaseURL + "/Ph-VaNhL")
	if err == nil {
		fmt.Printf("✅ Сервер доступен! Ответ: %s (Код: %d)\n", response, statusCode)
	} else {
		fmt.Printf("❌ Ошибка ping: %v\n", err)
	}

	// 2. Тестируем создание короткого URL через JSON (POST /api/shorten)
	fmt.Println("\n🔹 Тест: POST /api/shorten (JSON)")

	data, err := json.Marshal(rs)
	if err != nil {
		logger.Log.Fatalf("error: %s", err.Error())
	}
	response, statusCode, err = postClient.PostJSON(cfg.BaseURL+"/api/shorten", data)
	if err == nil {
		fmt.Printf("✅ Короткие URL созданы: %s (Код: %d)\n", response, statusCode)
	} else {
		fmt.Printf("❌ Ошибка создания JSON URL: %v\n", err)
	}

	// 3. Тестируем создание короткого URL через FormData (POST /)
	fmt.Println("\n🔹 Тест: POST / (FormData)")
	for _, r := range rs {
		response, statusCode, err = postClient.PostFormData(cfg.BaseURL+"/", r.OrigURL)
		if err == nil {
			fmt.Printf("✅ Короткий URL создан: %s (Код: %d)\n", response, statusCode)
		} else {
			fmt.Printf("❌ Ошибка создания FormData URL: %v\n", err)
		}
	}

	fmt.Println("\n✅ Все тесты завершены!")
}
