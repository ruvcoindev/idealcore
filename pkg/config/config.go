// Package config — конфигурация приложения
package config

import (
	"os"
	"strconv"
	"strings"
)

// Config — глобальная конфигурация
type Config struct {
	ServerPort        string
	ServerHost        string
	EnableYggdrasil   bool
	OllamaHost        string
	OllamaPort        string
	DefaultModel      string
	AlternativeModels []string
	VectorDim         int
	MaxConnections    int
	EFConstruction    int
	EFSearch          int
	RAGTopK           int
	HybridAlpha       float64
	DataDir           string
	TraumaDBPath      string
	DefenseDBPath     string
	YggdrasilAddr     string
	YggdrasilPeers    []string
	YggdrasilEnable   bool
}

// Load загружает конфигурацию
func Load() *Config {
	return &Config{
		ServerPort:        getEnv("SERVER_PORT", "8081"),
		ServerHost:        getEnv("SERVER_HOST", "localhost"),
		EnableYggdrasil:   getEnvBool("ENABLE_YGGDRASIL", false),
		OllamaHost:        getEnv("OLLAMA_HOST", "localhost"),
		OllamaPort:        getEnv("OLLAMA_PORT", "11434"),
		DefaultModel:      getEnv("DEFAULT_MODEL", "llama3.1:8b-instruct-q4_k_m"),
		AlternativeModels: getEnvSlice("ALTERNATIVE_MODELS", []string{"qwen3.5:9b-q4_k_m"}),
		VectorDim:         384,
		MaxConnections:    16,
		EFConstruction:    200,
		EFSearch:          50,
		RAGTopK:           5,
		HybridAlpha:       0.7,
		DataDir:           getEnv("DATA_DIR", "~/.idealcore/data"),
		TraumaDBPath:      getEnv("TRAUMA_DB", "~/.idealcore/traumas.json"),
		DefenseDBPath:     getEnv("DEFENSE_DB", "~/.idealcore/defenses.json"),
		YggdrasilAddr:     getEnv("YGGDRASIL_ADDR", "tcp://localhost:9001"),
		YggdrasilPeers:    getEnvSlice("YGGDRASIL_PEERS", []string{}),
		YggdrasilEnable:   getEnvBool("YGGDRASIL_ENABLE", false),
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

func getEnvSlice(key string, defaultVal []string) []string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
