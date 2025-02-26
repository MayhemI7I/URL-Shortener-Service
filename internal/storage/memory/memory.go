package memory

import (
   "context"
   "sync"

)

type MemoryStorage struct {
   urls map[string]string
   mu   sync.Mutex
}

func NewMemoryStorage() (*MemoryStorage, error) {
   return &MemoryStorage{
   	urls: make(map[string]string),
   	mu:   sync.Mutex{},
   }, nil
}

func (ms *MemoryStorage) Get(ctx context.Context, shortUrl string) (string, error) {
   select {
   case <-ctx.Done():
   	return "", ctx.Err()
   default:
   }
   ms.mu.Lock()
   defer ms.mu.Unlock()
   longUrl, ok := ms.urls[shortUrl]
   if !ok {
   	return "", nil
   }
   return longUrl, nil
}

func (ms *MemoryStorage) Save(ctx context.Context, shortUrl, longUrl string) error {
   select {
   case <-ctx.Done():
   	return ctx.Err()
   default:
   }
   ms.mu.Lock()
   defer ms.mu.Unlock()
   ms.urls[shortUrl] = longUrl
   return nil
}
func(ms *MemoryStorage) Close() error {
	return nil
}
func(ms *MemoryStorage) IfExistUrl(LongURL string)(string, error)  {
   ms.urls
}
