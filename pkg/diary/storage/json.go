// Package storage — JSON реализация хранилища
package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// JSONStorage — хранилище на основе JSON файлов
type JSONStorage struct {
	cfg     *Config
	dataDir string
	entries map[string]map[string]*Entry // userID -> questionID -> Entry
	mu      sync.RWMutex
	closed  bool
}

// init регистрирует JSON хранилище при загрузке пакета
func init() {
	Register("json", NewJSONStorage)
}

// NewJSONStorage создаёт новое JSON хранилище
func NewJSONStorage(cfg *Config) (Storage, error) {
	// Раскрываем ~ в пути
	dataDir := cfg.DataDir
	if len(dataDir) > 1 && dataDir[:2] == "~/" {
		home, _ := os.UserHomeDir()
		dataDir = filepath.Join(home, dataDir[2:])
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, &StorageError{
			Code:    "mkdir_failed",
			Message: "failed to create data directory",
			Err:     err,
		}
	}

	store := &JSONStorage{
		cfg:     cfg,
		dataDir: dataDir,
		entries: make(map[string]map[string]*Entry),
	}

	// Загружаем существующие данные
	if err := store.load(); err != nil {
		return nil, err
	}

	return store, nil
}

// Name возвращает имя реализации
func (s *JSONStorage) Name() string {
	return "json"
}

// Init инициализирует хранилище
func (s *JSONStorage) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return &StorageError{
			Code:    "closed",
			Message: "storage is closed",
		}
	}

	return nil
}

// Close закрывает хранилище
func (s *JSONStorage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Сохраняем все данные перед закрытием
	for userID := range s.entries {
		if err := s.save(userID); err != nil {
			return err
		}
	}

	s.closed = true
	return nil
}

// Save сохраняет запись
func (s *JSONStorage) Save(ctx context.Context, entry *Entry) error {
	if s.cfg.ReadOnly {
		return ErrReadOnly
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return &StorageError{Code: "closed", Message: "storage is closed"}
	}

	if entry == nil {
		return &StorageError{Code: "invalid_entry", Message: "entry is nil"}
	}

	if s.entries[entry.UserID] == nil {
		s.entries[entry.UserID] = make(map[string]*Entry)
	}

	questionKey := fmt.Sprintf("%s_%s", entry.Section, entry.ID)
	entry.UpdatedAt = time.Now()

	if _, exists := s.entries[entry.UserID][questionKey]; exists {
		entry.CreatedAt = s.entries[entry.UserID][questionKey].CreatedAt
	} else {
		entry.CreatedAt = time.Now()
	}

	s.entries[entry.UserID][questionKey] = entry
	return s.save(entry.UserID)
}

// Get получает запись
func (s *JSONStorage) Get(ctx context.Context, userID, questionID string) (*Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, &StorageError{Code: "closed", Message: "storage is closed"}
	}

	if userEntries, ok := s.entries[userID]; ok {
		for _, entry := range userEntries {
			if entry.ID == questionID || entry.Question == questionID {
				return entry, nil
			}
		}
	}

	return nil, ErrNotFound
}

// Update обновляет запись
func (s *JSONStorage) Update(ctx context.Context, entry *Entry) error {
	if s.cfg.ReadOnly {
		return ErrReadOnly
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return &StorageError{Code: "closed", Message: "storage is closed"}
	}

	if entry == nil {
		return &StorageError{Code: "invalid_entry", Message: "entry is nil"}
	}

	if s.entries[entry.UserID] == nil {
		return ErrNotFound
	}

	questionKey := fmt.Sprintf("%s_%s", entry.Section, entry.ID)
	if _, exists := s.entries[entry.UserID][questionKey]; !exists {
		return ErrNotFound
	}

	entry.UpdatedAt = time.Now()
	s.entries[entry.UserID][questionKey] = entry
	return s.save(entry.UserID)
}

// Delete удаляет запись
func (s *JSONStorage) Delete(ctx context.Context, userID, questionID string) error {
	if s.cfg.ReadOnly {
		return ErrReadOnly
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return &StorageError{Code: "closed", Message: "storage is closed"}
	}

	if userEntries, ok := s.entries[userID]; ok {
		for key := range userEntries {
			if userEntries[key].ID == questionID {
				delete(userEntries, key)
				return s.save(userID)
			}
		}
	}

	return ErrNotFound
}

// GetByUser получает все записи пользователя
func (s *JSONStorage) GetByUser(ctx context.Context, userID string) ([]*Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, &StorageError{Code: "closed", Message: "storage is closed"}
	}

	var result []*Entry
	if userEntries, ok := s.entries[userID]; ok {
		for _, entry := range userEntries {
			result = append(result, entry)
		}
	}

	return result, nil
}

// GetBySection получает записи по разделу
func (s *JSONStorage) GetBySection(ctx context.Context, userID, section string) ([]*Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, &StorageError{Code: "closed", Message: "storage is closed"}
	}

	var result []*Entry
	if userEntries, ok := s.entries[userID]; ok {
		for _, entry := range userEntries {
			if entry.Section == section {
				result = append(result, entry)
			}
		}
	}

	return result, nil
}

// GetByTag получает записи по тегу
func (s *JSONStorage) GetByTag(ctx context.Context, userID, tag string) ([]*Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, &StorageError{Code: "closed", Message: "storage is closed"}
	}

	var result []*Entry
	if userEntries, ok := s.entries[userID]; ok {
		for _, entry := range userEntries {
			for _, t := range entry.Tags {
				if t == tag {
					result = append(result, entry)
					break
				}
			}
		}
	}

	return result, nil
}

// Import импортирует записи
func (s *JSONStorage) Import(ctx context.Context, userID string, entries []*Entry) error {
	if s.cfg.ReadOnly {
		return ErrReadOnly
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return &StorageError{Code: "closed", Message: "storage is closed"}
	}

	if s.entries[userID] == nil {
		s.entries[userID] = make(map[string]*Entry)
	}

	for _, entry := range entries {
		questionKey := fmt.Sprintf("%s_%s", entry.Section, entry.ID)
		s.entries[userID][questionKey] = entry
	}

	return s.save(userID)
}

// Export экспортирует записи
func (s *JSONStorage) Export(ctx context.Context, userID string) ([]*Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, &StorageError{Code: "closed", Message: "storage is closed"}
	}

	var result []*Entry
	if userEntries, ok := s.entries[userID]; ok {
		for _, entry := range userEntries {
			result = append(result, entry)
		}
	}

	return result, nil
}

// Count возвращает количество записей
func (s *JSONStorage) Count(ctx context.Context, userID string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return 0, &StorageError{Code: "closed", Message: "storage is closed"}
	}

	if userEntries, ok := s.entries[userID]; ok {
		return len(userEntries), nil
	}

	return 0, nil
}

// GetStats возвращает статистику
func (s *JSONStorage) GetStats(ctx context.Context) (*Stats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, &StorageError{Code: "closed", Message: "storage is closed"}
	}

	stats := &Stats{
		SectionStats: make(map[string]int),
		TagStats:     make(map[string]int),
	}

	for _, userEntries := range s.entries {
		stats.TotalUsers++
		for _, entry := range userEntries {
			stats.TotalEntries++
			stats.SectionStats[entry.Section]++
			for _, tag := range entry.Tags {
				stats.TagStats[tag]++
			}
		}
	}

	// Размер хранилища
	files, _ := filepath.Glob(filepath.Join(s.dataDir, "*.json"))
	for _, file := range files {
		info, err := os.Stat(file)
		if err == nil {
			stats.StorageSize += info.Size()
		}
	}

	stats.LastUpdated = time.Now()
	return stats, nil
}

// load загружает данные с диска
func (s *JSONStorage) load() error {
	files, err := filepath.Glob(filepath.Join(s.dataDir, "*.json"))
	if err != nil {
		return err
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var entries map[string]*Entry
		if err := json.Unmarshal(data, &entries); err != nil {
			continue
		}

		userID := filepath.Base(file)
		userID = userID[:len(userID)-len(".json")]

		s.entries[userID] = entries
	}

	return nil
}

// save сохраняет данные пользователя на диск
func (s *JSONStorage) save(userID string) error {
	file := filepath.Join(s.dataDir, fmt.Sprintf("%s.json", userID))

	data, err := json.MarshalIndent(s.entries[userID], "", "  ")
	if err != nil {
		return &StorageError{
			Code:    "marshal_failed",
			Message: "failed to marshal entries",
			Err:     err,
		}
	}

	if err := os.WriteFile(file, data, 0600); err != nil {
		return &StorageError{
			Code:    "write_failed",
			Message: "failed to write file",
			Err:     err,
		}
	}

	return nil
}

// Проверка реализации интерфейса
var _ Storage = (*JSONStorage)(nil)
