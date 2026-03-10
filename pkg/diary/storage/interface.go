// Package storage предоставляет абстрактный интерфейс для хранилищ дневника
package storage

import (
	"context"
	"time"

	"github.com/ruvcoindev/idealcore/pkg/diary"
)

// Entry — запись дневника (дублируем из diary для независимости)
type Entry struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Section   string    `json:"section"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Storage — интерфейс для всех реализаций хранилищ
type Storage interface {
	// Инициализация
	Init() error
	Close() error

	// CRUD операции
	Save(ctx context.Context, entry *Entry) error
	Get(ctx context.Context, userID, questionID string) (*Entry, error)
	Update(ctx context.Context, entry *Entry) error
	Delete(ctx context.Context, userID, questionID string) error

	// Поиск
	GetByUser(ctx context.Context, userID string) ([]*Entry, error)
	GetBySection(ctx context.Context, userID, section string) ([]*Entry, error)
	GetByTag(ctx context.Context, userID, tag string) ([]*Entry, error)

	// Массовые операции
	Import(ctx context.Context, userID string, entries []*Entry) error
	Export(ctx context.Context, userID string) ([]*Entry, error)

	// Статистика
	Count(ctx context.Context, userID string) (int, error)
	GetStats(ctx context.Context) (*Stats, error)

	// Мета
	Name() string // имя реализации (json, sqlite, etc.)
}

// Stats — статистика хранилища
type Stats struct {
	TotalUsers   int            `json:"total_users"`
	TotalEntries int            `json:"total_entries"`
	SectionStats map[string]int `json:"section_stats"`
	TagStats     map[string]int `json:"tag_stats"`
	StorageSize  int64          `json:"storage_size_bytes"`
	LastUpdated  time.Time      `json:"last_updated"`
}

// Config — конфигурация хранилища
type Config struct {
	// Общие
	Type     string `json:"type"` // json, sqlite, postgres, etc.
	DataDir  string `json:"data_dir"`
	ReadOnly bool   `json:"read_only"`

	// SQLite
	SQLiteDSN string `json:"sqlite_dsn"`

	// PostgreSQL
	PostgresDSN string `json:"postgres_dsn"`

	// Кэширование
	CacheEnabled bool `json:"cache_enabled"`
	CacheTTL     int  `json:"cache_ttl_seconds"`

	// Шифрование
	EncryptionEnabled bool   `json:"encryption_enabled"`
	EncryptionKey     string `json:"encryption_key"`
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		Type:              "json",
		DataDir:           "~/.idealcore/diary",
		ReadOnly:          false,
		CacheEnabled:      true,
		CacheTTL:          300, // 5 минут
		EncryptionEnabled: false,
	}
}

// Validate проверяет конфигурацию
func (c *Config) Validate() error {
	if c.Type == "" {
		return ErrInvalidConfig
	}
	if c.DataDir == "" && c.Type == "json" {
		return ErrInvalidConfig
	}
	if c.Type == "sqlite" && c.SQLiteDSN == "" {
		return ErrInvalidConfig
	}
	if c.Type == "postgres" && c.PostgresDSN == "" {
		return ErrInvalidConfig
	}
	return nil
}

// Ошибки хранилища
var (
	ErrNotFound        = &StorageError{Code: "not_found", Message: "entry not found"}
	ErrAlreadyExists   = &StorageError{Code: "already_exists", Message: "entry already exists"}
	ErrInvalidConfig   = &StorageError{Code: "invalid_config", Message: "invalid storage configuration"}
	ErrReadOnly        = &StorageError{Code: "read_only", Message: "storage is read-only"}
	ErrEncryption      = &StorageError{Code: "encryption", Message: "encryption/decryption failed"}
	ErrConnection      = &StorageError{Code: "connection", Message: "database connection failed"}
	ErrTransaction     = &StorageError{Code: "transaction", Message: "transaction failed"}
)

// StorageError — ошибка хранилища
type StorageError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *StorageError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *StorageError) Unwrap() error {
	return e.Err
}

// Factory — фабрика для создания хранилищ
type Factory func(cfg *Config) (Storage, error)

// Registry — реестр зарегистрированных хранилищ
var Registry = make(map[string]Factory)

// Register регистрирует реализацию хранилища
func Register(name string, factory Factory) {
	Registry[name] = factory
}

// New создаёт хранилище по типу из конфига
func New(cfg *Config) (Storage, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	factory, ok := Registry[cfg.Type]
	if !ok {
		return nil, &StorageError{
			Code:    "unknown_type",
			Message: "unknown storage type: " + cfg.Type,
		}
	}

	return factory(cfg)
}
