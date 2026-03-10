// Package diary предоставляет хранилище для записей дневника
package diary

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry — запись дневника
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

// Store — хранилище записей дневника
type Store struct {
	dataDir string
	entries map[string]map[string]*Entry // userID -> questionID -> Entry
	mu      sync.RWMutex
}

// NewStore создаёт новое хранилище
func NewStore(dataDir string) (*Store, error) {
	if dataDir == "" {
		home, _ := os.UserHomeDir()
		dataDir = filepath.Join(home, ".idealcore", "diary")
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	store := &Store{
		dataDir: dataDir,
		entries: make(map[string]map[string]*Entry),
	}

	// Загружаем существующие данные
	if err := store.load(); err != nil {
		return nil, err
	}

	return store, nil
}

// SaveEntry сохраняет запись дневника
func (s *Store) SaveEntry(userID, section, answer string, tags []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if userID == "" {
		return fmt.Errorf("userID is required")
	}

	if section == "" {
		return fmt.Errorf("section is required")
	}

	// Получаем вопрос
	questions := GetQuestionsBySection(section)
	if len(questions) == 0 {
		return fmt.Errorf("no questions found for section: %s", section)
	}

	// Инициализируем мапу пользователя если нужно
	if s.entries[userID] == nil {
		s.entries[userID] = make(map[string]*Entry)
	}

	now := time.Now()

	// Создаём или обновляем запись
	for _, q := range questions {
		entryID := fmt.Sprintf("%s:%s:%s", userID, section, q.ID)
		
		if existing, ok := s.entries[userID][q.ID]; ok {
			// Обновляем существующую
			existing.Answer = answer
			existing.Tags = tags
			existing.UpdatedAt = now
			s.entries[userID][q.ID] = existing
		} else {
			// Создаём новую
			s.entries[userID][q.ID] = &Entry{
				ID:        entryID,
				UserID:    userID,
				Section:   section,
				Question:  q.Text,
				Answer:    answer,
				Tags:      tags,
				CreatedAt: now,
				UpdatedAt: now,
			}
		}
	}

	// Сохраняем на диск
	return s.save(userID)
}

// GetEntries возвращает все записи пользователя
func (s *Store) GetEntries(userID string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]string)
	if userEntries, ok := s.entries[userID]; ok {
		for questionID, entry := range userEntries {
			result[questionID] = entry.Answer
		}
	}

	return result, nil
}

// GetEntry возвращает конкретную запись
func (s *Store) GetEntry(userID, questionID string) (*Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if userEntries, ok := s.entries[userID]; ok {
		entry, exists := userEntries[questionID]
		return entry, exists
	}

	return nil, false
}

// DeleteEntry удаляет запись
func (s *Store) DeleteEntry(userID, questionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if userEntries, ok := s.entries[userID]; ok {
		delete(userEntries, questionID)
		return s.save(userID)
	}

	return fmt.Errorf("entry not found")
}

// GetAllEntries возвращает все записи всех пользователей (для админа)
func (s *Store) GetAllEntries() map[string]map[string]*Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]map[string]*Entry)
	for userID, entries := range s.entries {
		result[userID] = make(map[string]*Entry)
		for qID, entry := range entries {
			result[userID][qID] = entry
		}
	}
	return result
}

// GetEntriesBySection возвращает записи пользователя по разделу
func (s *Store) GetEntriesBySection(userID, section string) ([]*Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

// GetEntriesByTag возвращает записи по тегу
func (s *Store) GetEntriesByTag(userID, tag string) ([]*Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

// load загружает данные с диска
func (s *Store) load() error {
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

		// Извлекаем userID из имени файла
		userID := filepath.Base(file)
		userID = userID[:len(userID)-len(".json")]

		s.entries[userID] = entries
	}

	return nil
}

// save сохраняет данные пользователя на диск
func (s *Store) save(userID string) error {
	file := filepath.Join(s.dataDir, fmt.Sprintf("%s.json", userID))

	data, err := json.MarshalIndent(s.entries[userID], "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0600)
}

// Export exports all entries for a user to a single JSON file
func (s *Store) Export(userID string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, ok := s.entries[userID]
	if !ok {
		return nil, fmt.Errorf("no entries found for user: %s", userID)
	}

	return json.MarshalIndent(entries, "", "  ")
}

// Import imports entries from a JSON file
func (s *Store) Import(userID string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var entries map[string]*Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}

	s.entries[userID] = entries
	return s.save(userID)
}

// Count возвращает количество записей пользователя
func (s *Store) Count(userID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if userEntries, ok := s.entries[userID]; ok {
		return len(userEntries)
	}
	return 0
}

// GetStats возвращает статистику по всем пользователям
func (s *Store) GetStats() StoreStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalUsers := len(s.entries)
	totalEntries := 0
	sectionCounts := make(map[string]int)

	for _, userEntries := range s.entries {
		totalEntries += len(userEntries)
		for _, entry := range userEntries {
			sectionCounts[entry.Section]++
		}
	}

	return StoreStats{
		TotalUsers:   totalUsers,
		TotalEntries: totalEntries,
		SectionStats: sectionCounts,
	}
}

// StoreStats — статистика хранилища
type StoreStats struct {
	TotalUsers   int            `json:"total_users"`
	TotalEntries int            `json:"total_entries"`
	SectionStats map[string]int `json:"section_stats"`
}
