// Package vector — pure-Go векторное хранилище
package vector

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

// VectorStore — интерфейс для векторных хранилищ
type VectorStore interface {
	Add(id string, vector []float32, metadata map[string]interface{}) error
	Search(query []float32, k int) ([]SearchResult, error)
	Remove(id string) bool
	Count() int
}

// VectorItem — элемент с вектором
type VectorItem struct {
	ID       string
	Metadata map[string]interface{}
	Vector   []float32
}

// Store — реализация VectorStore
type Store struct {
	mu    sync.RWMutex
	items map[string]VectorItem
	dim   int
}

// NewStore создаёт новое хранилище
func NewStore(dim int) *Store {
	return &Store{
		items: make(map[string]VectorItem),
		dim:   dim,
	}
}

// Add добавляет вектор
func (s *Store) Add(id string, vector []float32, metadata map[string]interface{}) error {
	if len(vector) != s.dim {
		return fmt.Errorf("dimension mismatch: expected %d, got %d", s.dim, len(vector))
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[id] = VectorItem{ID: id, Vector: vector, Metadata: metadata}
	return nil
}

// SearchResult — результат поиска
type SearchResult struct {
	ID         string
	Similarity float32
	Metadata   map[string]interface{}
}

// Search находит ближайшие векторы
func (s *Store) Search(query []float32, k int) ([]SearchResult, error) {
	if len(query) != s.dim {
		return nil, fmt.Errorf("query dimension mismatch")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []SearchResult
	qn := norm(query)
	for _, item := range s.items {
		if len(item.Vector) != s.dim {
			continue
		}
		sim := cosineSimilarity(query, item.Vector, qn, norm(item.Vector))
		results = append(results, SearchResult{
			ID:         item.ID,
			Similarity: sim,
			Metadata:   item.Metadata,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})
	if len(results) > k {
		results = results[:k]
	}
	return results, nil
}

// Remove удаляет вектор
func (s *Store) Remove(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.items[id]
	if ok {
		delete(s.items, id)
	}
	return ok
}

// Count возвращает количество элементов
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// cosineSimilarity вычисляет косинусное сходство
func cosineSimilarity(a, b []float32, na, nb float32) float32 {
	if na == 0 || nb == 0 {
		return 0
	}
	dot := float32(0)
	for i := range a {
		dot += a[i] * b[i]
	}
	return dot / (na * nb)
}

// norm вычисляет L2-норму
func norm(v []float32) float32 {
	sum := float32(0)
	for _, x := range v {
		sum += x * x
	}
	return float32(math.Sqrt(float64(sum)))
}

// Проверка реализации интерфейса
var _ VectorStore = (*Store)(nil)
