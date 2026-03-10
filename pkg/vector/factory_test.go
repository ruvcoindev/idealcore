package vector

import "testing"

func TestFactory_PureGo(t *testing.T) {
	cfg := Config{
		Backend:   "purego",
		Dimension: 3,
	}
	
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	
	// Проверяем, что работает
	_ = store.Add("test", []float32{1, 0, 0}, nil)
	results, err := store.Search([]float32{1, 0, 0}, 1)
	if err != nil {
		t.Errorf("search failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestFactory_HNSW_Fallback(t *testing.T) {
	// Если CGO не доступен — должен сфоллбэкиться на Pure-Go
	cfg := Config{
		Backend:      "hnsw",
		Dimension:    3,
		MaxElements:  1000,
		M:            16,
		EfConstruction: 200,
	}
	
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("factory failed: %v", err)
	}
	
	// Должен работать в любом случае
	_ = store.Add("fallback_test", []float32{0, 1, 0}, nil)
	if store.Count() != 1 {
		t.Error("fallback store not working")
	}
}
