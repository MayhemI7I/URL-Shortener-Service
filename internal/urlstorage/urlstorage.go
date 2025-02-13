package urlstorage

import (
	"encoding/json"
	"errors"
	"local/logger"
	"log"
	"os"
	"sync"
)





type URLStorage struct {
	urls map[string]string 
	mu   sync.Mutex
	file *os.File
}

func NewURLStorage(filename string) (*URLStorage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	storage := &URLStorage{
		urls: make(map[string]string),
		mu:   sync.Mutex{},
		file: file,
  
	}
	return storage, nil
}

func (us *URLStorage) Close() error {
	return us.file.Close()
	
}

func (us *URLStorage) Load() error {
	us.mu.Lock()
	defer us.mu.Unlock()
	decoder := json.NewDecoder(us.file)
	for {
		var entry map[string]string
		if err := decoder.Decode(&entry); err != nil {
			break
		}
		for key, value := range entry {
			us.urls[key] = value
		}
	}
	return nil
}

func (us *URLStorage) SaveURL(shortURL, longURL string) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	if shortURL == "" || longURL == "" {
		log.Printf("Invalid argument: %s, %s", shortURL, longURL)
		return errors.New("invalid argument")
	}
	if _, exists := us.urls[shortURL]; exists {
		log.Printf("URL already exists: %s", shortURL)
		return errors.New("URL already exists")
	}
	us.urls[shortURL] = longURL

	encoder := json.NewEncoder(us.file)
	if err := encoder.Encode(us.urls); err != nil {
		return err
	}


	logger.Log.Info("Saved: %s -> %s", shortURL, longURL)
	return nil
}


func (us *URLStorage) GetURL(shortUrl string) (string, error) {
	us.mu.Lock()
	defer us.mu.Unlock()
	if shortUrl == "" {
		log.Printf("Invalid argument: %s", shortUrl)
		return "", errors.New("invalid short URL argument")
  
	}
	value, ok := us.urls[shortUrl]
	if !ok {
		return "", errors.New("URL not found in storage")
	}
	logger.Log.Info("Retrived: %s -> %s", shortUrl, value)
	return value, nil

}

