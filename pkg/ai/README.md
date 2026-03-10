
# pkg/ai — Интеграция с Ollama

Модуль для работы с локальным LLM-сервером [Ollama](https://ollama.ai) в проекте idealcore.

## Возможности

- ✅ Генерация намерений для практики на гвоздях с психологическим контекстом
- ✅ Получение эмбеддингов для векторного поиска (RAG)
- ✅ Анализ записей дневника на предмет травм, защит, чакр
- ✅ Поддержка нескольких моделей с авто-выбором под задачу
- ✅ Пост-обработка ответов: очистка от артефактов генерации
- ✅ RAG-поиск по фрагментам дневника для персонализации намерений

## Быстрый старт

### 1. Установи модели в Ollama

```bash
# Базовая модель (английский + код)
ollama pull llama3.1:8b-instruct-q4_k_m

# Для лучшего русского языка (рекомендуется)
ollama pull qwen3.5:9b-q4_k_m

# Для эмбеддингов (векторный поиск)
ollama pull all-minilm:latest

# Альтернатива для эмбеддингов (лучше качество, больше размер)
ollama pull nomic-embed-text:latest

### 2. Запусти Ollama-сервер

# В отдельном терминале
ollama serve

# Или как системный сервис (Ubuntu)
sudo systemctl start ollama
sudo systemctl enable ollama

### 3. Проверь, что модели доступны

curl http://localhost:11434/api/tags

Ответ должен содержать список установленных моделей.

### 4. Используй в коде idealcore

package main

import (
    "context"
    "github.com/ruvcoindev/idealcore/pkg/ai"
    "github.com/ruvcoindev/idealcore/pkg/config"
    "github.com/ruvcoindev/idealcore/pkg/psychology"
    "github.com/ruvcoindev/idealcore/pkg/vector"
)

func main() {
    cfg := config.Load()
    traumaDB := psychology.NewTraumaDB()
    vectorStore := vector.NewStore(384)
    
    client := ai.NewOllamaClient(cfg, traumaDB, vectorStore)
    
    ctx := context.Background()
    
    // Генерация намерения
    intention, err := client.GenerateIntention(ctx, ai.IntentionContext{
        UserID:      "user-123",
        Theme:       "границы с родителями",
        TraumaTypes: []domain.TraumaType{domain.TraumaNarcissisticParent},
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(intention)
}

#### Стиль генерации намерений

Модель настроена на выдачу текста в следующем стиле:

Принципы

Принцип	                                 Описание
Твёрдость	        Без эзотерической «ваты», слащавости, пустых обещаний
Настоящее время	        Местоимение «Я», глаголы в настоящем («я есть», не «я буду»)
Структура	        Заземление → Право быть → Отпускание контроля → Принятие силы → Финал
Теневые аспекты  	Работа с виной, страхом, контролем, обязанностью — не избегание
Ритм   	                Короткие рубленые фразы, чередующиеся с плавными переходами
Тон	                Принятие без снисхождения, твёрдость без агрессии


#### Пример вывода

Я здесь. В этом теле. В этом моменте.

Я имею право занять пространство — не заслуживая, не извиняясь.
Мои ноги чувствуют опору. Моя спина чувствует поддержку.

Я отпускаю необходимость быть удобной для других.
Я отпускаю вину за то, что выбираю себя.

Моя сила — не в том, чтобы выдержать всё.
А в том, чтобы выбрать, что достойно меня.

Боль, которую я чувствую — это не наказание.
Это сигнал: ты живая. Ты здесь. Ты возвращаешься домой.

Я есть. Этого достаточно.


#### Настройка
###  Выбор модели
Модель выбирается автоматически по теме намерения:

Тип                                   темы	                               Модель	                          Причина
Русскоязычные (`отношения`, `границы`, `вина`, `стыд`, `родители`)	`qwen3.5:9b-q4_k_m`	Лучшее понимание нюансов русского языка
Технические (`код`, `анализ данных`, `структура`)	                `llama3.1:8b-q4_k_m`	Быстрее, лучше для логики
Эмбеддинги	                                                         `all-minilm:latest`    Оптимизирована для векторного поиска


##### Параметры генерации
Настройки по умолчанию в OllamaOptions:

Options: OllamaOptions{
    Temperature: 0.3,    // Низкая для точности, не креатива
    TopP:        0.9,    // Разнообразие без потери фокуса
    NumPredict:  512,    // Макс. длина ответа (токены)
    Stop:        []string{"\n\n###", "USER:", "ASSISTANT:"}, // Маркеры остановки
}

#### Изменение параметров
Через конфигурацию в pkg/config/config.go:

type Config struct {
    OllamaHost          string
    OllamaPort          string
    DefaultModel        string
    AlternativeModels   []string
    // ... другие поля
}

Или через переменные окружения:

export OLLAMA_HOST=localhost
export OLLAMA_PORT=11434
export DEFAULT_MODEL=llama3.1:8b-instruct-q4_k_m
export ALTERNATIVE_MODELS=qwen3.5:9b-q4_k_m,llama3.2:3b-q4_k_m


#### RAG-контекст (поиск по дневнику)
Для персонализации намерений можно добавить релевантные фрагменты из дневника пользователя:

// 1. Поиск фрагментов по теме
snippets, err := client.SearchDiarySnippets(ctx, "чувство вины", userID, k=3)
if err != nil {
    log.Fatal(err)
}

// 2. Передача в контекст генерации
ic := ai.IntentionContext{
    UserID:        userID,
    Theme:         "право на ошибку",
    DiarySnippets: snippets, // RAG-контекст
}

// 3. Генерация с учётом истории
intention, err := client.GenerateIntention(ctx, ic)

### Как это работает

1. Запрос пользователя векторизуется через GetEmbedding()
2. Векторный поиск находит похожие записи дневника (косинусное сходство)
3. Фрагменты добавляются в промпт как контекст
4. Модель генерирует намерение с учётом личной истории пользователя

##### Анализ записей дневника

 Модуль может автоматически анализировать текст дневника и извлекать психологические маркеры:

entryText := "Я чувствую вину, когда отдыхаю. Кажется, что если я не работаю — я бесполезен."

trauma, defenses, chakras, err := client.AnalyzeDiaryEntry(ctx, entryText)
if err != nil {
    log.Fatal(err)
}

// trauma: domain.TraumaType (например, "narcissistic_parent")
// defenses: []domain.DefenseMechanism (например, ["rationalization", "sublimation"])
// chakras: []domain.ChakraType (например, ["solar", "heart"])



Распознаваемые типы травм

ptsd — Посттравматическое стрессовое расстройство
cptsd — Комплексное ПТСР
bpd — Пограничное расстройство личности
npd — Нарциссическое расстройство личности
attachment — Травма привязанности
narcissistic_parent — Травма от нарциссического родителя
abandonment — Травма брошенности
betrayal — Травма предательства
shame — Травма стыда

Распознаваемые защитные механизмы

denial — Отрицание
projection — Проекция
rationalization — Рационализация
displacement — Смещение
sublimation — Сублимация
regression — Регрессия
reaction_formation — Реактивное образование
intellectualization — Интеллектуализация
narcissistic_grandiosity — Нарциссическая грандиозность
narcissistic_devaluation — Нарциссическое обесценивание
triangulation — Триангуляция

Распознаваемые чакры

root — Муладхара (безопасность, выживание)
sacral — Свадхистана (эмоции, творчество)
solar — Манипура (воля, границы)
heart — Анахата (любовь, принятие)
throat — Вишудха (выражение, правда)
third_eye — Аджна (интуиция, видение)
crown — Сахасрара (связь, трансценденция)

Тестирование
Юнит-тесты

##bash

# Запустить все тесты
go test ./pkg/ai/... -v

# Запустить конкретный тест
go test ./pkg/ai/... -v -run TestBuildSystemPrompt

# Запустить с покрытием
go test ./pkg/ai/... -v -cover

Интеграционные тесты
Требуют запущенного Ollama-сервера:

## bash

# Убедись, что Ollama работает
curl http://localhost:11434/api/tags

# Запусти интеграционный тест
go test ./pkg/ai/... -v -run Integration

Моки для тестирования без Ollama
Для тестирования без реального сервера создай мок:

#go

type MockOllamaClient struct {
    GenerateFunc func(ctx context.Context, ic IntentionContext) (string, error)
}

func (m *MockOllamaClient) GenerateIntention(ctx context.Context, ic IntentionContext) (string, error) {
    return m.GenerateFunc(ctx, ic)
}


Troubleshooting
«connection refused»
Проблема: Ollama-сервер не запущен.
Решение:

### bash 

# Проверь статус
ollama serve

# Или перезапусти сервис
sudo systemctl restart ollama

Медленная генерация (менее 5 токенов/сек)
Проблема: Нехватка RAM или тяжелая модель.
Решение:

### bash 

# Используй квантованные модели
ollama pull llama3.1:8b-instruct-q4_k_m

# Закрой тяжёлые приложения (браузер, IDE)
# Освободи минимум 8 GB RAM

# Для 16 GB RAM оптимально: q4_k_m или q5_k_m
# Для 8 GB RAM: q3_k_m или q4_k_s


Модель выдаёт «пластиковый» текст
Проблема: Слишком общая генерация, клише.
Решение: Добавь уточнение в RawPrompt:

####go

ic := ai.IntentionContext{
    Theme: "право на отдых",
    RawPrompt: "Сделай текст более телесным, убери общие фразы типа 'вселенная' и 'поток'. Добавь конкретики про дыхание, ноги, опору.",
}


Ошибка парсинга JSON в AnalyzeDiaryEntry
Проблема: Модель добавляет маркеры кода (```json) в ответ.
Решение: Функция postProcess() уже очищает ответ. Если проблема сохраняется — уменьши температуру:


##go

Options: OllamaOptions{
    Temperature: 0.1, // Минимум креатива для JSON
}

Эмбеддинги не совпадают по размерности
Проблема: Разные модели эмбеддингов дают разную размерность.
Решение: Используй одну модель для всех эмбеддингов:

# all-minilm даёт 384 измерения
ollama pull all-minilm:latest

# В конфиге укажи:
export EMBEDDING_MODEL=all-minilm:latest
export VECTOR_DIM=384

Производительность
Бенчмарки на i7-9700, 16 GB RAM

Модель	                 Размер	  Скорость генерации	RAM	Рекомендация
`llama3.1:8b-q4_k_m`	~4.7 GB	  8-12 ток/сек	       6 GB	✅ Базовая
`qwen3.5:9b-q4_k_m`	~5.5 GB	  6-9 ток/сек	       7 GB	✅ Для русского
`llama3.2:3b-q4_k_m`	~2.0 GB	  20-30 ток/сек	       3 GB	✅ Быстрая, но проще
`all-minilm:latest`	~0.5 GB	  50+ ток/сек	       1 GB	✅ Для эмбеддингов

Оптимизация

## bash

# 1. Используй GPU если есть (NVIDIA)
# Ollama автоматически использует CUDA при наличии

# 2. Для Intel GPU (как у тебя UHD 630)
# Поддержка ограничена, но можно попробовать:
export OLLAMA_NUM_GPU=0  # force CPU

# 3. Увеличь контекст если нужно
export OLLAMA_CONTEXT_LENGTH=8192  # по умолчанию 2048

Архитектура

pkg/ai/
├── ollama.go           # Основной клиент
├── ollama_test.go      # Тесты
├── README.md           # Документация
├── prompts.go          # Шаблоны промптов (опционально)
└── cache.go            # Кэш ответов (опционально)

Зависимости

### go

require (
    github.com/ruvcoindev/idealcore/internal/domain
    github.com/ruvcoindev/idealcore/pkg/config
    github.com/ruvcoindev/idealcore/pkg/psychology
    github.com/ruvcoindev/idealcore/pkg/vector
)

Лицензия

Apache 2.0 — как и весь проект idealcore.

Авторы

Виталий — архитектура, психологический контекст
idealcore contributors — развитие модуля

Связанные модули

pkg/vector — Векторное хранилище для RAG
pkg/diary — Дневник «Кто я»
pkg/psychology — База травм и защит
pkg/rag — RAG-система с гибридным поиском

