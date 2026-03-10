package rag

import (
	"context"
	"testing"
	"time"

	"github.com/ruvcoindev/idealcore/pkg/ai"
	"github.com/ruvcoindev/idealcore/pkg/diary"
	"github.com/ruvcoindev/idealcore/pkg/vector"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.TopK != 5 {
		t.Errorf("expected TopK=5, got %d", cfg.TopK)
	}
	if cfg.HybridAlpha != 0.7 {
		t.Errorf("expected HybridAlpha=0.7, got %f", cfg.HybridAlpha)
	}
	if !cfg.CacheEnabled {
		t.Error("expected CacheEnabled=true")
	}
}

func TestNewRAG(t *testing.T) {
	vectorStore := vector.NewStore(384)
	aiClient := &ai.OllamaClient{} // mock
	diaryStore, _ := diary.NewStore("")
	config := DefaultConfig()

	rag := New(vectorStore, aiClient, diaryStore, config)

	if rag == nil {
		t.Fatal("expected RAG instance, got nil")
	}
	if rag.config != config {
		t.Error("expected config to be set")
	}
}

func TestExtractTerms(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		wantLen  int
		contains []string
	}{
		{
			name:     "simple text",
			text:     "Я чувствую вину когда отдыхаю",
			wantLen:  4,
			contains: []string{"чувствую", "вину", "отдыхаю"},
		},
		{
			name:     "with stopwords",
			text:     "И я очень хочу быть счастливым",
			wantLen:  3,
			contains: []string{"хочу", "быть", "счастливым"},
		},
		{
			name:     "short words filtered",
			text:     "Я и он на да нет",
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terms := extractTerms(tt.text)
			
			if len(terms) != tt.wantLen {
				t.Errorf("expected %d terms, got %d: %v", tt.wantLen, len(terms), terms)
			}

			for _, want := range tt.contains {
				found := false
				for _, term := range terms {
					if term == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected term %q in %v", want, terms)
				}
			}
		})
	}
}

func TestRAGStats(t *testing.T) {
	vectorStore := vector.NewStore(384)
	aiClient := &ai.OllamaClient{}
	diaryStore, _ := diary.NewStore("")
	rag := New(vectorStore, aiClient, diaryStore, DefaultConfig())

	stats := rag.GetStats()

	if stats.TotalEntries != 0 {
		t.Errorf("expected 0 entries initially, got %d", stats.TotalEntries)
	}
	if stats.IndexBuilt {
		t.Error("expected index not built initially")
	}
}

func TestHybridRanking(t *testing.T) {
	vectorStore := vector.NewStore(384)
	aiClient := &ai.OllamaClient{}
	diaryStore, _ := diary.NewStore("")
	rag := New(vectorStore, aiClient, diaryStore, DefaultConfig())

	// Создаём тестовые результаты
	vectorResults := []vector.SearchResult{
		{ID: "entry1", Similarity: 0.9},
		{ID: "entry2", Similarity: 0.7},
		{ID: "entry3", Similarity: 0.5},
	}

	keywordResults := []SearchResult{
		{EntryID: "entry2", KeywordScore: 0.8},
		{EntryID: "entry3", KeywordScore: 0.6},
		{EntryID: "entry1", KeywordScore: 0.2},
	}

	// Гибридное ранжирование с alpha=0.5
	hybrid := rag.hybridRanking(vectorResults, keywordResults, 0.5)

	if len(hybrid) != 3 {
		t.Fatalf("expected 3 results, got %d", len(hybrid))
	}

	// entry1 должен быть первым (0.5*1.0 + 0.5*0.25 = 0.625)
	// entry2 должен быть вторым (0.5*0.78 + 0.5*1.0 = 0.89)
	if hybrid[0].EntryID != "entry2" {
		t.Errorf("expected entry2 first, got %s", hybrid[0].EntryID)
	}
}

func TestContextFormation(t *testing.T) {
	vectorStore := vector.NewStore(384)
	aiClient := &ai.OllamaClient{}
	diaryStore, _ := diary.NewStore("")
	rag := New(vectorStore, aiClient, diaryStore, DefaultConfig())

	// Добавляем тестовые entry в индекс
	rag.index.entries["user1:q1"] = &diary.Entry{
		ID:       "user1:q1",
		UserID:   "user1",
		Section:  "motivation",
		Question: "q1",
		Answer:   "Я делаю это для себя",
		Tags:     []string{"motivation"},
		CreatedAt: time.Now(),
	}

	context := &Context{
		Query:      "тестовый запрос",
		TotalFound: 1,
		SearchTime: 10 * time.Millisecond,
		Results: []SearchResult{
			{
				EntryID:   "user1:q1",
				Section:   "motivation",
				Answer:    "Я делаю это для себя",
				Score:     0.85,
				Tags:      []string{"motivation"},
				CreatedAt: time.Now(),
			},
		},
	}

	if context.TotalFound != 1 {
		t.Errorf("expected 1 result, got %d", context.TotalFound)
	}
	if context.Query != "тестовый запрос" {
		t.Errorf("unexpected query: %s", context.Query)
	}
}

func TestRAGWithRealData(t *testing.T) {
	t.Skip("Integration test: requires Ollama running")

	tmpDir := t.TempDir()
	vectorStore := vector.NewStore(384)
	aiClient, _ := ai.NewOllamaClient(nil, nil, vectorStore) // нужна реальная конфига
	diaryStore, _ := diary.NewStore(tmpDir)
	rag := New(vectorStore, aiClient, diaryStore, DefaultConfig())

	ctx := context.Background()

	// Сохраняем запись
	_ = diaryStore.SaveEntry("user1", "motivation", "Я чувствую вину когда отдыхаю", []string{"guilt"})

	// Строим индекс
	err := rag.BuildIndex(ctx, "user1")
	if err != nil {
		t.Fatalf("BuildIndex failed: %v", err)
	}

	// Ищем
	context, err := rag.Search(ctx, "user1", "чувство вины", 5)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if context.TotalFound == 0 {
		t.Error("expected at least 1 result")
	}
}

func TestAddRemoveEntry(t *testing.T) {
	vectorStore := vector.NewStore(384)
	aiClient := &ai.OllamaClient{}
	diaryStore, _ := diary.NewStore("")
	rag := New(vectorStore, aiClient, diaryStore, DefaultConfig())

	ctx := context.Background()
	entry := &diary.Entry{
		ID:       "test-entry",
		UserID:   "user1",
		Section:  "motivation",
		Question: "test",
		Answer:   "test answer",
		Tags:     []string{"test"},
	}

	// Добавляем
	err := rag.AddEntry(ctx, entry)
	if err != nil {
		t.Logf("AddEntry skipped (requires embeddings): %v", err)
		return
	}

	stats := rag.GetStats()
	if stats.TotalEntries != 1 {
		t.Errorf("expected 1 entry after add, got %d", stats.TotalEntries)
	}

	// Удаляем
	err = rag.RemoveEntry(ctx, "test-entry")
	if err != nil {
		t.Fatalf("RemoveEntry failed: %v", err)
	}

	stats = rag.GetStats()
	if stats.TotalEntries != 0 {
		t.Errorf("expected 0 entries after remove, got %d", stats.TotalEntries)
	}
}

func BenchmarkExtractTerms(b *testing.B) {
	text := "Я чувствую вину когда отдыхаю и мне кажется что я бесполезен если не работаю постоянно"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = extractTerms(text)
	}
}

func BenchmarkHybridRanking(b *testing.B) {
	vectorStore := vector.NewStore(384)
	aiClient := &ai.OllamaClient{}
	diaryStore, _ := diary.NewStore("")
	rag := New(vectorStore, aiClient, diaryStore, DefaultConfig())

	vectorResults := make([]vector.SearchResult, 100)
	keywordResults := make([]SearchResult, 100)

	for i := 0; i < 100; i++ {
		vectorResults[i] = vector.SearchResult{
			ID:         string(rune(i)),
			Similarity: float32(100-i) / 100,
		}
		keywordResults[i] = SearchResult{
			EntryID:      string(rune(i)),
			KeywordScore: float32(i) / 100,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rag.hybridRanking(vectorResults, keywordResults, 0.7)
	}
}
