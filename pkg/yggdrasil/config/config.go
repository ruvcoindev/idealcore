// Package config предоставляет конфигурацию для Yggdrasil-транспорта
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// YggdrasilConfig — конфигурация Yggdrasil
type YggdrasilConfig struct {
	// Сетевые настройки
	ListenAddresses []string `json:"listen_addresses"`
	Peers           []string `json:"peers"`
	Multicast       bool     `json:"multicast"`

	// Криптография
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`

	// Синхронизация
	SyncInterval    int  `json:"sync_interval"`     // секунды
	EnableSync      bool `json:"enable_sync"`
	ConflictStrategy string `json:"conflict_strategy"`

	// Логирование
	LogLevel string `json:"log_level"`
	LogFile  string `json:"log_file"`

	// Производительность
	MaxConnections   int `json:"max_connections"`
	MessageQueueSize int `json:"message_queue_size"`
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *YggdrasilConfig {
	return &YggdrasilConfig{
		ListenAddresses:  []string{"tcp://:9999"},
		Peers:            []string{},
		Multicast:        true,
		SyncInterval:     300, // 5 минут
		EnableSync:       true,
		ConflictStrategy: "last_write_wins",
		LogLevel:         "info",
		LogFile:          "",
		MaxConnections:   50,
		MessageQueueSize: 1000,
	}
}

// Load загружает конфигурацию из файла
func Load(path string) (*YggdrasilConfig, error) {
	cfg := DefaultConfig()

	// Путь по умолчанию
	if path == "" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".idealcore", "yggdrasil.json")
	}

	// Читаем файл если существует
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Создаём файл с дефолтной конфигурацией
			os.MkdirAll(filepath.Dir(path), 0755)
			saveConfig(path, cfg)
			return cfg, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save сохраняет конфигурацию в файл
func (c *YggdrasilConfig) Save(path string) error {
	return saveConfig(path, c)
}

func saveConfig(path string, cfg *YggdrasilConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// GetYggdrasilSocketPath возвращает путь к сокету Yggdrasil
func GetYggdrasilSocketPath() string {
	// Linux
	if _, err := os.Stat("/var/run/yggdrasil.sock"); err == nil {
		return "unix:///var/run/yggdrasil.sock"
	}
	// macOS
	if _, err := os.Stat("/tmp/yggdrasil.sock"); err == nil {
		return "unix:///tmp/yggdrasil.sock"
	}
	// Fallback на TCP
	return "tcp://localhost:9001"
}

// Validate проверяет конфигурацию
func (c *YggdrasilConfig) Validate() error {
	if c.SyncInterval < 60 {
		return fmt.Errorf("sync_interval must be at least 60 seconds")
	}
	if c.MaxConnections < 1 {
		return fmt.Errorf("max_connections must be at least 1")
	}
	if c.MessageQueueSize < 100 {
		return fmt.Errorf("message_queue_size must be at least 100")
	}
	return nil
}
