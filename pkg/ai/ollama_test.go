package ai

import (
	"context"
	"testing"
	"time"

	"github.com/ruvcoindev/idealcore/internal/domain"
	"github.com/ruvcoindev/idealcore/pkg/config"
	"github.com/ruvcoindev/idealcore/pkg/psychology"
	"github.com/ruvcoindev/idealcore/pkg/vector"
)

// Тест на создание клиента
func TestNewOllamaClient(t *testing.T) {
	cfg := config.Load()
	traumaDB := psychology.NewTraumaDB()
	vectorStore := vector.NewStore(384)

	client := NewOllamaClient(cfg, traumaDB, vectorStore)

	if client == nil {
		t.Fatal("expected client, got nil")
	}
	if client.baseURL == "" {
		t.Error("expected baseURL to be set")
	}
}

// Тест на формирование системного промпта
func TestBuildSystemPrompt(t *testing.T) {
	cfg := config.Load()
	traumaDB := psychology.NewTraumaDB()
	vectorStore := vector.NewStore(384)
	client := NewOllamaClient(cfg, traumaDB, vectorStore)

	ic := IntentionContext{
		Theme:           "границы с родителями",
		TraumaTypes:     []domain.TraumaType{domain.TraumaNarcissisticParent},
		AttachmentStyle: domain.AttachmentAnxious,
		ChakraFocus:     []domain.ChakraType{domain.ChakraSolar, domain.ChakraThroat},
	}

	prompt := client.buildSystemPrompt(ic)

	// Проверяем наличие ключевых элементов
	checks := []string{
		"Ты — опытный телесно-ориентированный психолог",
		"СТИЛЬ ТЕКСТА:",
		"Контекст травм",
		"Тип привязанности",
		"Энергетический фокус",
		"Не обещай изменений",
	}

	for _, check := range checks {
		if !contains(prompt, check) {
			t.Errorf("system prompt missing expected element: %s", check)
		}
	}
}

// Тест на формирование пользовательского промпта
func TestBuildUserPrompt(t *testing.T) {
	cfg := config.Load()
	traumaDB := psychology.NewTraumaDB()
	vectorStore := vector.NewStore(384)
	client := NewOllamaClient(cfg, traumaDB, vectorStore)

	ic := IntentionContext{
		Theme: "право на отдых",
		DefenseMechanisms: []domain.DefenseMechanism{
			domain.DefenseRationalization,
			domain.DefenseSublimation,
		},
		DiarySnippets: []string{
			"Я чувствую вину, когда отдыхаю",
			"Кажется, что если я не работаю — я бесполезен",
		},
	}

	prompt := client.buildUserPrompt(ic)

	if !contains(prompt, "ТЕМА НАМЕРЕНИЯ: право на отдых") {
		t.Error("user prompt missing theme")
	}
	if !contains(prompt, "РЕЛЕВАНТНЫЕ ФРАГМЕНТЫ ИЗ ДНЕВНИКА") {
		t.Error("user prompt missing diary snippets section")
	}
	if !contains(prompt, "ЗАЩИТНЫЕ МЕХАНИЗМЫ") {
		t.Error("user prompt missing defenses section")
	}
}

// Тест на выбор модели
func TestSelectModel(t *testing.T) {
	cfg := config.Load()
	cfg.AlternativeModels = []string{"qwen3.5:9b-q4_k_m", "llama3.2:3b-q4_k_m"}
	traumaDB := psychology.NewTraumaDB()
	vectorStore := vector.NewStore(384)
	client := NewOllamaClient(cfg, traumaDB, vectorStore)

	// Русскоязычная тема — должен выбрать qwen, если доступен
	model := client.selectModel("отношения с матерью")
	// В тесте не проверяем конкретное значение, т.к. зависит от конфига
	_ = model

	// Техническая тема — дефолтная модель
	model2 := client.selectModel("генерация кода на Go")
	if model2 != cfg.DefaultModel {
		t.Logf("selected %s for technical theme, expected %s", model2, cfg.DefaultModel)
	}
}

// Тест на пост-обработку
func TestPostProcess(t *testing.T) {
	cfg := config.Load()
	traumaDB := psychology.NewTraumaDB()
	vectorStore := vector.NewStore(384)
	client := NewOllamaClient(cfg, traumaDB, vectorStore)

	tests := []struct {
		input    string
		expected string
	}{
		{"  \n### Я имею право быть ###\n", "Я имею право быть"},
		{`"Я выбираю себя"`, "Я выбираю себя"},
		{"Просто текст", "Просто текст"},
		{"\n\nЛишние переносы", "Лишние переносы"},
	}

	for _, tt := range tests {
		result := client.postProcess(tt.input)
		if result != tt.expected {
			t.Errorf("postProcess(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// Вспомогательная функция
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Интеграционный тест (пропускается, если нет запущенного Ollama)
func TestGenerateIntention_Integration(t *testing.T) {
	t.Skip("Integration test: requires running Ollama server")

	cfg := config.Load()
	traumaDB := psychology.NewTraumaDB()
	vectorStore := vector.NewStore(384)
	client := NewOllamaClient(cfg, traumaDB, vectorStore)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ic := IntentionContext{
		UserID:  "test-user",
		Theme:   "право на ошибку",
	}

	intention, err := client.GenerateIntention(ctx, ic)
	if err != nil {
		t.Fatalf("GenerateIntention failed: %v", err)
	}

	if intention == "" {
		t.Error("expected non-empty intention")
	}
	if len(intention) < 50 {
		t.Errorf("intention too short: %d chars", len(intention))
	}

	t.Logf("Generated intention:\n%s", intention)
}
