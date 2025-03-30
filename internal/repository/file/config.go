package file

	import (
	"os"

	"fmt"
	"github.com/spf13/pflag"
)

type FileStorageConfig struct {
	name string
	flags *pflag.FlagSet
	path string
}
// Создает новую конфигурацию файлового хранилища
func NewFileStorageConfig() *FileStorageConfig {
	return &FileStorageConfig{
		name: "FileStorage",
		flags: pflag.NewFlagSet("FileStorage", pflag.ExitOnError),
		path: "./data/urls.json",
	}
}
// Добавляет флаги для конфигурации файлового хранилища
func (c *FileStorageConfig) AddFlags() {
	c.flags.StringVar(&c.path, "file-storage-path", c.path, "Path to the file storage")
}

// Загружает конфигурацию из переменных окружения
func (c *FileStorageConfig) LoadFromEnv() {
	c.path = os.Getenv("FILE_STORAGE_PATH")
}

// Проверяет корректность конфигурации
func (c *FileStorageConfig) Validate() error {
	if c.path == "" {
		return fmt.Errorf("path to the file storage is not set")
	}
	return nil
}

// Возвращает путь к файлу хранилища
func (c *FileStorageConfig) GetPath() string {
	return c.path
}

// Возвращает флаги конфигурации
func (c *FileStorageConfig) GetFlags() *pflag.FlagSet {
	return c.flags
}

