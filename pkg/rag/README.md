# pkg/rag — RAG (Retrieval-Augmented Generation) система

Модуль для семантического поиска и извлечения контекста из дневника для генерации намерений.

## Возможности

- ✅ Гибридный поиск (векторный + keyword)
- ✅ Автоматическая индексация записей дневника
- ✅ Reranking результатов через LLM (опционально)
- ✅ Кэширование результатов поиска
- ✅ Извлечение релевантного контекста для AI
- ✅ Поддержка тегов и разделов
- ✅ Статистика и метрики поиска

## Быстрый старт

### 1. Создание RAG системы

go
package main

import (
    "context"
    "github.com/ruvcoindev/idealcore/pkg/rag"
    "github.com/ruvcoindev/idealcore/pkg/vector"
    "github.com/ruvcoindev/idealcore/pkg/ai"
    "github.com/ruvcoindev/idealcore/pkg/diary"
)

func main() {
    vectorStore := vector.NewStore(384)
    aiClient, _ := ai.NewOllamaClient(cfg, traumaDB, vectorStore)
    diaryStore, _ := diary.NewStore("~/.idealcore/diary")
    
    config := rag.DefaultConfig()
    config.TopK = 5
    config.HybridAlpha = 0.7 // 70% векторный, 30% keyword
    
    ragSystem := rag.New(vectorStore, aiClient, diaryStore, config)
    
    // Строим индекс
    ctx := context.Background()
    _ = ragSystem.BuildIndex(ctx, "user-123")
    
    // Ищем
    context, _ := ragSystem.Search(ctx, "user-123", "чувство вины", 5)
    
    for _, res := range context.Results {
        println(res.Answer)
    }
}

### 2. Поиск с контекстом

```go

// Извлекаем контекст для генерации намерения
contextText, err := ragSystem.GetContext(ctx, "user-123", "границы с родителями")

// contextText содержит:
// РЕЛЕВАНТНЫЕ ФРАГМЕНТЫ ИЗ ДНЕВНИКА:
// 
// 1. [boundaries] Я чувствую вину когда говорю нет
//    Теги: guilt, boundaries
// 
// 2. [motivation] Я делаю это чтобы меня любили
//    Теги: approval, love


Конфигурация
Параметры поиска

Параметр	                       По умолчанию	            Описание
TopK	                                    5	           Количество результатов
HybridAlpha	                           0.7	           Вес векторного поиска (0..1)
MinSimilarity	                           0.3	           Минимальное сходство
MaxContextLength	                   2048            Макс. длина контекста (токены)


Кэширование

Параметр	По умолчанию	   Описание
CacheEnabled	true	         Включить кэш
CacheTTL	300	         Время жизни кэша (сек)

Reranking

Параметр	По умолчанию	Описание
RerankEnabled	true	      Включить reranking
RerankTopK	10	      Сколько результатов rerankить

Фрагментация

Параметр	    По умолчанию	       Описание
ChunkSize	      500	           Размер чанка (символы)
ChunkOverlap	      50	           Перекрытие чанков

API
Основные методы

Метод	                                                                 Описание
New(vectorStore, aiClient, diaryStore, config)	                  Создать RAG систему
BuildIndex(ctx, userID)	                                          Построить индекс для пользователя
Search(ctx, userID, query, topK)	                          Выполнить поиск
GetContext(ctx, userID, theme)	                                  Извлечь контекст для генерации
AddEntry(ctx, entry)	                                          Добавить запись в индекс
RemoveEntry(ctx, entryID)	                                  Удалить запись из индекса
GetStats()	                                                  Получить статистику

SearchResult

```go

type SearchResult struct {
    EntryID       string
    UserID        string
    Section       string
    Question      string
    Answer        string
    Score         float32    // гибридный score
    VectorScore   float32    // векторная часть
    KeywordScore  float32    // keyword часть
    Tags          []string
    CreatedAt     time.Time
}


Алгоритм работы
1. Индексация

Запись дневника
    ↓
Токенизация → Keywords индекс
    ↓
Эмбеддинг → Векторный индекс
    ↓
Сохранение в хранилище

2. Поиск

Запрос пользователя
    ↓
    ├─→ Векторный поиск (косинусное сходство)
    │       ↓
    │   Vector scores
    │
    └─→ Keyword поиск (TF-IDF упрощённый)
            ↓
        Keyword scores
            ↓
    Гибридное ранжирование (alpha * vector + (1-alpha) * keyword)
            ↓
        Reranking (LLM, опционально)
            ↓
        Фильтрация по min_similarity
            ↓
        TopK результатов

Гибридный поиск
Формула

Score = alpha * VectorScore + (1 - alpha) * KeywordScore

где:
alpha = 0.7 (по умолчанию)
VectorScore — косинусное сходство (0..1)
KeywordScore — доля совпавших терминов (0..1)
Пример

Запрос: "чувство вины когда отдыхаю"

Результаты:
1. Entry A: Vector=0.9, Keyword=0.8 → Score=0.87
2. Entry B: Vector=0.8, Keyword=0.6 → Score=0.74
3. Entry C: Vector=0.6, Keyword=0.9 → Score=0.69

Reranking
Reranking использует LLM для более точной оценки релевантности:

```go

// Запрос к LLM
"Ранжируй следующие фрагменты по релевантности запросу: \"чувство вины\"

1. Я чувствую вину когда отдыхаю
2. Мне стыдно за свои ошибки
3. Я боюсь осуждения других

Верни номера в порядке убывания релевантности: "

// Ответ LLM: "1,2,3"

Интеграция с другими модулями

Модуль	Интеграция
pkg/ai  	Генерация эмбеддингов, reranking
pkg/diary	Источник данных
pkg/vector	Векторное хранилище
pkg/web  	HTTP API для поиска

Примеры использования
Поиск похожих записей

```go

// Пользователь пишет запись
answer := "Я чувствую вину когда отдыхаю"

// Ищем похожие
context, _ := rag.Search(ctx, userID, answer, 5)

// Показываем паттерны
for _, res := range context.Results {
    fmt.Printf("[%s] %s (score: %.2f)\n", 
        res.Section, res.Answer, res.Score)
}

Генерация намерения с контекстом

```go

// Извлекаем контекст
contextText, _ := rag.GetContext(ctx, userID, "границы")

// Передаём в AI
intention, _ := aiClient.GenerateIntention(ctx, ai.IntentionContext{
    Theme: "границы",
    DiaryContext: contextText, // RAG контекст
})

Анализ паттернов

```go

// Ищем все записи про вину
context, _ := rag.Search(ctx, userID, "вина", 100)

// Считаем по разделам
sectionCounts := make(map[string]int)
for _, res := range context.Results {
    sectionCounts[res.Section]++
}

// Выводим
for section, count := range sectionCounts {
    fmt.Printf("%s: %d записей\n", section, count)
}

Производительность
Бенчмарки (i7-9700, 16GB RAM)

Операция	100 записей	1000 записей	10000 записей
BuildIndex	  ~1 сек	~10 сек	         ~100 сек
Search	          ~50 мс	~100 мс	         ~300 мс
AddEntry	  ~20 мс	~20 мс	         ~20 мс


Оптимизация
Кэширование: Включите CacheEnabled для повторяющихся запросов
Инкрементальная индексация: Используйте AddEntry вместо BuildIndex
Ограничьте TopK: Не запрашивайте больше чем нужно
Отключите reranking: Если скорость важнее точности
Troubleshooting
«no results found»
Проблема: Поиск не находит результаты.
Решение:

```go

// Уменьшите min_similarity
config.MinSimilarity = 0.1

// Увеличьте topK
context, _ := rag.Search(ctx, userID, query, 20)

«slow search»
Проблема: Поиск работает медленно.
Решение:

```go

// Отключите reranking
config.RerankEnabled = false

// Уменьшите HybridAlpha (больше keyword, быстрее)
config.HybridAlpha = 0.3

«index out of sync»
Проблема: Индекс не соответствует данным.
Решение:

```go 

// Перестройте индекс
_ = rag.BuildIndex(ctx, userID)

// Или используйте инкрементальное обновление
_ = rag.AddEntry(ctx, newEntry)
_ = rag.RemoveEntry(ctx, oldEntryID)

Тестирование

```bash

# Юнит-тесты
go test ./pkg/rag/... -v

# Бенчмарки
go test ./pkg/rag/... -bench=. -benchmem

# Integration test (требует Ollama)
go test ./pkg/rag/... -run Integration -v

Лицензия
Apache 2.0 — как и весь проект idealcore.




