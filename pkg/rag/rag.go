// Package rag предоставляет RAG (Retrieval-Augmented Generation) систему
// для семантического поиска и извлечения контекста из дневника
package rag

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ruvcoindev/idealcore/pkg/ai"
	"github.com/ruvcoindev/idealcore/pkg/diary"
	"github.com/ruvcoindev/idealcore/pkg/vector"
)

// RAG — Retrieval-Augmented Generation система
type RAG struct {
	vectorStore vector.VectorStore
	aiClient    *ai.OllamaClient
	diaryStore  *diary.Store
	config      *Config
	index       *Index
	mu          sync.RWMutex
}

// Config — конфигурация RAG системы
type Config struct {
	// Поиск
	TopK              int     `json:"top_k"`                // количество результатов
	HybridAlpha       float32 `json:"hybrid_alpha"`         // вес векторного поиска (0..1)
	MinSimilarity     float32 `json:"min_similarity"`       // минимальное сходство
	MaxContextLength  int     `json:"max_context_length"`   // макс. длина контекста (токены)
	
	// Кэширование
	CacheEnabled      bool    `json:"cache_enabled"`
	CacheTTL          int     `json:"cache_ttl_seconds"`
	
	// Reranking
	RerankEnabled     bool    `json:"rerank_enabled"`
	RerankTopK        int     `json:"rerank_top_k"`
	
	// Фрагментация
	ChunkSize         int     `json:"chunk_size"`           // размер чанка (символы)
	ChunkOverlap      int     `json:"chunk_overlap"`        // перекрытие чанков
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		TopK:             5,
		HybridAlpha:      0.7,      // 70% векторный, 30% keyword
		MinSimilarity:    0.3,
		MaxContextLength: 2048,
		CacheEnabled:     true,
		CacheTTL:         300,      // 5 минут
		RerankEnabled:    true,
		RerankTopK:       10,
		ChunkSize:        500,
		ChunkOverlap:     50,
	}
}

// SearchResult — результат поиска
type SearchResult struct {
	EntryID    string    `json:"entry_id"`
	UserID     string    `json:"user_id"`
	Section    string    `json:"section"`
	Question   string    `json:"question"`
	Answer     string    `json:"answer"`
	Score      float32   `json:"score"`
	VectorScore float32  `json:"vector_score"`
	KeywordScore float32 `json:"keyword_score"`
	Tags       []string  `json:"tags"`
	CreatedAt  time.Time `json:"created_at"`
}

// Context — извлечённый контекст для генерации
type Context struct {
	Query      string         `json:"query"`
	Results    []SearchResult `json:"results"`
	TotalFound int            `json:"total_found"`
	SearchTime time.Duration  `json:"search_time"`
}

// Index — индекс для быстрого поиска
type Index struct {
	entries   map[string]*diary.Entry
	vectors   map[string][]float32
	keywords  map[string][]string // term -> entryIDs
	built     bool
	builtAt   time.Time
}

// New создаёт новую RAG систему
func New(vectorStore vector.VectorStore, aiClient *ai.OllamaClient, diaryStore *diary.Store, config *Config) *RAG {
	if config == nil {
		config = DefaultConfig()
	}

	rag := &RAG{
		vectorStore: vectorStore,
		aiClient:    aiClient,
		diaryStore:  diaryStore,
		config:      config,
		index: &Index{
			entries:  make(map[string]*diary.Entry),
			vectors:  make(map[string][]float32),
			keywords: make(map[string][]string),
		},
	}

	return rag
}

// BuildIndex строит индекс для всех записей пользователя
func (r *RAG) BuildIndex(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	entries, err := r.diaryStore.GetEntries(userID)
	if err != nil {
		return fmt.Errorf("failed to get entries: %w", err)
	}

	// Очищаем старый индекс
	r.index = &Index{
		entries:  make(map[string]*diary.Entry),
		vectors:  make(map[string][]float32),
		keywords: make(map[string][]string),
	}

	// Индексируем каждую запись
	for questionID, answer := range entries {
		entry := &diary.Entry{
			ID:       fmt.Sprintf("%s:%s", userID, questionID),
			UserID:   userID,
			Question: questionID,
			Answer:   answer,
		}

		// Получаем эмбеддинг
		embedding, err := r.aiClient.GetEmbedding(ctx, answer, "all-minilm:latest")
		if err != nil {
			continue
		}

		// Сохраняем в индекс
		r.index.entries[entry.ID] = entry
		r.index.vectors[entry.ID] = embedding

		// Извлекаем ключевые слова
		terms := extractTerms(answer)
		for _, term := range terms {
			r.index.keywords[term] = append(r.index.keywords[term], entry.ID)
		}

		// Добавляем в векторное хранилище
		_ = r.vectorStore.Add(entry.ID, embedding, map[string]interface{}{
			"user_id":  userID,
			"question": questionID,
			"answer":   answer,
		})
	}

	r.index.built = true
	r.index.builtAt = time.Now()

	return nil
}

// Search выполняет гибридный поиск (vector + keyword)
func (r *RAG) Search(ctx context.Context, userID, query string, topK int) (*Context, error) {
	startTime := time.Now()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.index.built {
		if err := r.BuildIndex(ctx, userID); err != nil {
			return nil, err
		}
	}

	if topK == 0 {
		topK = r.config.TopK
	}

	// 1. Векторный поиск
	queryEmbedding, err := r.aiClient.GetEmbedding(ctx, query, "all-minilm:latest")
	if err != nil {
		return nil, fmt.Errorf("failed to get query embedding: %w", err)
	}

	vectorResults, err := r.vectorStore.Search(queryEmbedding, topK*r.config.RerankTopK)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 2. Keyword поиск
	keywordResults := r.keywordSearch(query, userID)

	// 3. Гибридное ранжирование
	hybridResults := r.hybridRanking(vectorResults, keywordResults, r.config.HybridAlpha)

	// 4. Reranking (опционально)
	if r.config.RerankEnabled {
		hybridResults = r.rerank(ctx, query, hybridResults)
	}

	// 5. Фильтрация по минимальному сходству
	var filtered []SearchResult
	for _, res := range hybridResults {
		if res.Score >= r.config.MinSimilarity {
			filtered = append(filtered, res)
		}
	}

	// 6. Ограничиваем topK
	if len(filtered) > topK {
		filtered = filtered[:topK]
	}

	searchTime := time.Since(startTime)

	return &Context{
		Query:      query,
		Results:    filtered,
		TotalFound: len(filtered),
		SearchTime: searchTime,
	}, nil
}

// GetContext извлекает контекст для генерации намерения
func (r *RAG) GetContext(ctx context.Context, userID, theme string) (string, error) {
	context, err := r.Search(ctx, userID, theme, r.config.TopK)
	if err != nil {
		return "", err
	}

	if len(context.Results) == 0 {
		return "", nil
	}

	// Формируем текстовый контекст
	var sb strings.Builder
	sb.WriteString("РЕЛЕВАНТНЫЕ ФРАГМЕНТЫ ИЗ ДНЕВНИКА:\n\n")

	for i, res := range context.Results {
		sb.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, res.Section, res.Answer))
		if len(res.Tags) > 0 {
			sb.WriteString(fmt.Sprintf("   Теги: %s\n", strings.Join(res.Tags, ", ")))
		}
		sb.WriteString("\n")

		// Проверяем лимит длины
		if len(sb.String()) > r.config.MaxContextLength*4 { // примерная оценка символов
			break
		}
	}

	return sb.String(), nil
}

// keywordSearch выполняет поиск по ключевым словам
func (r *RAG) keywordSearch(query, userID string) []SearchResult {
	terms := extractTerms(query)
	entryScores := make(map[string]float32)

	// Подсчитываем совпадения терминов
	for _, term := range terms {
		if entryIDs, ok := r.index.keywords[term]; ok {
			for _, entryID := range entryIDs {
				if strings.HasPrefix(entryID, userID) {
					entryScores[entryID]++
				}
			}
		}
	}

	// Нормализуем scores и конвертируем в результаты
	var results []SearchResult
	for entryID, score := range entryScores {
		if entry, ok := r.index.entries[entryID]; ok {
			results = append(results, SearchResult{
				EntryID:      entryID,
				UserID:       entry.UserID,
				Section:      entry.Section,
				Question:     entry.Question,
				Answer:       entry.Answer,
				KeywordScore: score / float32(len(terms)),
				Tags:         entry.Tags,
				CreatedAt:    entry.CreatedAt,
			})
		}
	}

	// Сортируем по score
	sort.Slice(results, func(i, j int) bool {
		return results[i].KeywordScore > results[j].KeywordScore
	})

	return results
}

// hybridRanking объединяет векторный и keyword поиск
func (r *RAG) hybridRanking(vectorResults []vector.SearchResult, keywordResults []SearchResult, alpha float32) []SearchResult {
	// Создаём мапу для быстрого доступа
	keywordMap := make(map[string]float32)
	for _, res := range keywordResults {
		keywordMap[res.EntryID] = res.KeywordScore
	}

	// Нормализуем векторные scores
	maxVectorScore := float32(0)
	for _, res := range vectorResults {
		if res.Similarity > maxVectorScore {
			maxVectorScore = res.Similarity
		}
	}

	// Объединяем результаты
	combined := make(map[string]*SearchResult)

	for _, res := range vectorResults {
		entryID := res.ID
		var entry *diary.Entry
		
		// Ищем entry в индексе
		for id, e := range r.index.entries {
			if id == entryID {
				entry = e
				break
			}
		}

		if entry == nil {
			continue
		}

		vectorScore := res.Similarity / maxVectorScore
		keywordScore := keywordMap[entryID]

		hybridScore := alpha*vectorScore + (1-alpha)*keywordScore

		combined[entryID] = &SearchResult{
			EntryID:      entryID,
			UserID:       entry.UserID,
			Section:      entry.Section,
			Question:     entry.Question,
			Answer:       entry.Answer,
			VectorScore:  vectorScore,
			KeywordScore: keywordScore,
			Score:        hybridScore,
			Tags:         entry.Tags,
			CreatedAt:    entry.CreatedAt,
		}
	}

	// Конвертируем в слайс
	var results []SearchResult
	for _, res := range combined {
		results = append(results, *res)
	}

	// Сортируем по hybrid score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// rerank выполняет reranking с помощью LLM
func (r *RAG) rerank(ctx context.Context, query string, results []SearchResult) []SearchResult {
	if len(results) == 0 {
		return results
	}

	// Берём top K для reranking
	topK := r.config.RerankTopK
	if len(results) < topK {
		topK = len(results)
	}

	// Формируем запрос к LLM для reranking
	var sb strings.Builder
	sb.WriteString("Ранжируй следующие фрагменты по релевантности запросу: \"")
	sb.WriteString(query)
	sb.WriteString("\"\n\n")

	for i, res := range results[:topK] {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, res.Answer))
	}

	sb.WriteString("\nВерни номера фрагментов в порядке убывания релевантности (например: 3,1,2)")

	// TODO: Вызов LLM для reranking
	// Пока просто возвращаем как есть

	return results[:topK]
}

// AddEntry добавляет запись в индекс
func (r *RAG) AddEntry(ctx context.Context, entry *diary.Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Получаем эмбеддинг
	embedding, err := r.aiClient.GetEmbedding(ctx, entry.Answer, "all-minilm:latest")
	if err != nil {
		return err
	}

	// Добавляем в индекс
	r.index.entries[entry.ID] = entry
	r.index.vectors[entry.ID] = embedding

	// Извлекаем термины
	terms := extractTerms(entry.Answer)
	for _, term := range terms {
		r.index.keywords[term] = append(r.index.keywords[term], entry.ID)
	}

	// Добавляем в векторное хранилище
	return r.vectorStore.Add(entry.ID, embedding, map[string]interface{}{
		"user_id":  entry.UserID,
		"section":  entry.Section,
		"question": entry.Question,
		"answer":   entry.Answer,
		"tags":     entry.Tags,
	})
}

// RemoveEntry удаляет запись из индекса
func (r *RAG) RemoveEntry(ctx context.Context, entryID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Удаляем из индекса
	delete(r.index.entries, entryID)
	delete(r.index.vectors, entryID)

	// Удаляем из keyword индекса
	for term, entryIDs := range r.index.keywords {
		var newIDs []string
		for _, id := range entryIDs {
			if id != entryID {
				newIDs = append(newIDs, id)
			}
		}
		if len(newIDs) == 0 {
			delete(r.index.keywords, term)
		} else {
			r.index.keywords[term] = newIDs
		}
	}

	// Удаляем из векторного хранилища
	r.vectorStore.Remove(entryID)

	return nil
}

// GetStats возвращает статистику RAG системы
func (r *RAG) GetStats() RAGStats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return RAGStats{
		TotalEntries:  len(r.index.entries),
		TotalTerms:    len(r.index.keywords),
		IndexBuilt:    r.index.built,
		IndexBuiltAt:  r.index.builtAt,
		VectorCount:   r.vectorStore.Count(),
	}
}

// RAGStats — статистика RAG системы
type RAGStats struct {
	TotalEntries  int       `json:"total_entries"`
	TotalTerms    int       `json:"total_terms"`
	IndexBuilt    bool      `json:"index_built"`
	IndexBuiltAt  time.Time `json:"index_built_at"`
	VectorCount   int       `json:"vector_count"`
}

// extractTerms извлекает термины из текста для keyword поиска
func extractTerms(text string) []string {
	// Простая токенизация
	text = strings.ToLower(text)
	text = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'а' && r <= 'я' || r >= '0' && r <= '9' {
			return r
		}
		return ' '
	}, text)

	words := strings.Fields(text)

	// Фильтруем стоп-слова
	stopWords := map[string]bool{
		"и": true, "в": true, "во": true, "не": true, "что": true, "он": true,
		"на": true, "я": true, "с": true, "со": true, "как": true, "а": true,
		"то": true, "по": true, "только": true, "се": true, "эта": true,
		"мы": true, "бы": true, "же": true, "вы": true, "ты": true, "ну": true,
		"очень": true, "после": true, "даже": true, "нет": true, "так": true,
		"или": true, "при": true, "о": true, "из": true, "за": true, "до": true,
	}

	var terms []string
	for _, word := range words {
		if len(word) > 2 && !stopWords[word] {
			terms = append(terms, word)
		}
	}

	return terms
}
