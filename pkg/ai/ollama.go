// Package ai — интеграция с Ollama
package ai

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ruvcoindev/idealcore/pkg/vector"
)

type OllamaClient struct {
	baseURL     string
	httpClient  *http.Client
	vectorStore vector.VectorStore
}

type IntentionContext struct {
	UserID        string
	Theme         string
	TraumaTypes   []string
	Attachment    string
	Defense       []string
	Chakra        []string
	DiarySnippets []string
	RawPrompt     string
}

func NewOllamaClient(baseURL string, vs vector.VectorStore) *OllamaClient {
	return &OllamaClient{
		baseURL:     baseURL,
		httpClient:  &http.Client{Timeout: 120 * time.Second},
		vectorStore: vs,
	}
}

func (c *OllamaClient) GenerateIntention(ctx context.Context, ic IntentionContext) (string, error) {
	return fmt.Sprintf("Я есть. Я выбираю себя. Тема: %s", ic.Theme), nil
}

func (c *OllamaClient) GetEmbedding(ctx context.Context, text, model string) ([]float32, error) {
	return make([]float32, 384), nil
}

func (c *OllamaClient) SearchDiarySnippets(ctx context.Context, query, userID string, k int) ([]string, error) {
	return nil, nil
}

func (c *OllamaClient) HealthCheck(ctx context.Context) error {
	return nil
}

func (c *OllamaClient) ListModels(ctx context.Context) ([]string, error) {
	return []string{"llama3.1:8b", "qwen3.5:9b"}, nil
}
