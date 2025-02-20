package main

import (
	"bufio"
	"fmt"
	"local/config"
	"local/logger"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// Интерфейс для клиента, выполняющего запросы
// type PostClient interface {
// 	PostJSON(url string, longURL string) (string, int, error)
// 	PostFormData(url string, longURL string) (string, int, error)
// }

// Структура для реализации клиента с resty
type ClientReq struct {
	request *resty.Client
}

// Реализация метода PostJSON для отправки JSON
func (c *ClientReq) PostJSON(url, longURL string) (string, int, error) {
	response, err := c.request.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"url": longURL}).
		Post(url)

	if err != nil {
		logger.Log.Errorf("error: %s", err.Error())
		return "", 0, err
	}

	if response.StatusCode() != 200 && response.StatusCode() != 201 {
		logger.Log.Errorf("Server returned error: %s, %v", response.Status(), response.String())
		return "", response.StatusCode(), err
	}

	result := response.String()
	fmt.Println(response)
	return result, response.StatusCode(), nil
}

// Реализация метода PostFormData для отправки данных формы
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
		return "", response.StatusCode(), err
	}

	result := response.String()
	return result, response.StatusCode(), nil
}

func(c *ClientReq) GetPing(url string) (string, int, error) {
	response, err := c.request.R().
	Get(url)
		
	if response.StatusCode() != 200 && response.StatusCode() != 201 {
		logger.Log.Errorf("Server returned error: %s, %v", response.Status(), response.String())
		return "", response.StatusCode(), err
	}
	

	result := response.String()
	return result, response.StatusCode(), nil
}

// Чтение длинного URL с консоли
func readLongURL() (string, error) {
	fmt.Println("Введите длинный URL")
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

	logger.Log.Info("Starting server")

	// Чтение длинного URL
	long, err := readLongURL()
	if err != nil {
		logger.Log.Fatalf("error: %s", err.Error())
	}

	// Создание клиента для выполнения запросов
	client := resty.New()
	client.SetTimeout(500 * time.Millisecond)
	postClient := &ClientReq{request: client}

	var response string
	var statusCode int

	// В зависимости от BaseURL, выбираем способ отправки данных
	switch cfg.BaseURL {
	case "http://localhost:8080":
		response, statusCode, err = postClient.PostJSON(cfg.BaseURL, long)
		if err != nil || response == "" {
			logger.Log.Fatalf("Error posting JSON: %s", err.Error())
		}
	case "http://localhost:8080/api/shorten":
		response, statusCode, err = postClient.PostFormData(cfg.BaseURL, long)
		if err != nil || response == "" {
			logger.Log.Fatalf("Error posting form data: %s", err.Error())
		}
	case "http://localhost:8080/ping":
		response, statusCode, err = postClient.GetPing(cfg.BaseURL)
		if err != nil {
			logger.Log.Errorf("Error during GET request: %v", err)
		}
		fmt.Printf("Response body: %s\n", response)
		fmt.Printf("Status Code: %d\n", statusCode)
	}

	// Вывод результатов
	fmt.Println("Response:", response)
	fmt.Println("Status Code:", statusCode)
}
