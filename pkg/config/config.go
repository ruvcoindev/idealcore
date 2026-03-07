// Package config содержит конфигурацию приложения idealcore
package config

import (
	"os"
	"strconv"
)

// Config — глобальная конфигурация
type Config struct {
	// Сервер
	ServerPort     string
	ServerHost     string
	EnableYggdrasil bool

	// Ollama
	OllamaHost     string
	OllamaPort     string
	DefaultModel   string
	AlternativeModels []string

	// Векторное хранилище
	VectorDim      int
	MaxConnections int
	EFConstruction int
	EFSearch       int

	// RAG
	RAGTopK        int
	HybridAlpha    float64 // вес векторного поиска (0..1)

	// Пользователи
	DataDir        string

	// Психология
	TraumaDBPath   string
	DefenseDBPath  string
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	return &Config{
		ServerPort:        getEnv("SERVER_PORT", "8081"), // Не 8080!
		ServerHost:        getEnv("SERVER_HOST", "localhost"),
		EnableYggdrasil:   getEnvBool("ENABLE_YGGDRASIL", false),

		OllamaHost:        getEnv("OLLAMA_HOST", "localhost"),
		OllamaPort:        getEnv("OLLAMA_PORT", "11434"),
		DefaultModel:      getEnv("DEFAULT_MODEL", "llama3.1:8b-instruct-q4_k_m"),
		AlternativeModels: []string{"qwen3.5:9b-q4_k_m", "llama3.2:3b-q4_k_m"},

		VectorDim:         768, // для all-MiniLM-L6-v2
		MaxConnections:    16,
		EFConstruction:    200,
		EFSearch:          50,

		RAGTopK:           5,
		HybridAlpha:       0.7,

		DataDir:           getEnv("DATA_DIR", "./data"),
		TraumaDBPath:      getEnv("TRAUMA_DB", "./data/traumas.json"),
		DefenseDBPath:     getEnv("DEFENSE_DB", "./data/defenses.json"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}
