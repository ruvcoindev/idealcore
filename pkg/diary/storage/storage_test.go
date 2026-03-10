package storage

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Type != "json" {
		t.Errorf("expected type json, got %s", cfg.Type)
	}
	if cfg.DataDir != "~/.idealcore/diary" {
		t.Errorf("unexpected dataDir: %s", cfg.DataDir)
	}
	if !cfg.CacheEnabled {
		t.Error("expected cache enabled by default")
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{"valid json", &Config{Type: "json", DataDir: "/tmp/test"}, false},
		{"empty type", &Config{Type: "", DataDir: "/tmp"}, true},
		{"json no dir", &Config{Type: "json", DataDir: ""}, true},
		{"sqlite no dsn", &Config{Type: "sqlite", SQLiteDSN: ""}, true},
		{"postgres no dsn", &Config{Type: "postgres", PostgresDSN: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewJSONStorage(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{
		Type:    "json",
		DataDir: tmpDir,
	}

	store, err := NewJSONStorage(cfg)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer store.Close()

	if store.Name() != "json" {
		t.Errorf("expected name json, got %s", store.Name())
	}
}

func TestJSONStorageCRUD(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{
		Type:    "json",
		DataDir: tmpDir,
	}

	store, _ := NewJSONStorage(cfg)
	defer store.Close()

	ctx := context.Background()
	entry := &Entry{
		ID:       "test-1",
		UserID:   "user-123",
		Section:  "motivation",
		Question: "Ради кого я это делаю?",
		Answer:   "Для себя",
		Tags:     []string{"motivation", "self"},
	}

	// Save
	err := store.Save(ctx, entry)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Get
	got, err := store.Get(ctx, "user-123", "test-1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Answer != entry.Answer {
		t.Errorf("expected answer %s, got %s", entry.Answer, got.Answer)
	}

	// Update
	entry.Answer = "Обновлённый ответ"
	err = store.Update(ctx, entry)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, _ = store.Get(ctx, "user-123", "test-1")
	if got.Answer != "Обновлённый ответ" {
		t.Errorf("expected updated answer, got %s", got.Answer)
	}

	// Delete
	err = store.Delete(ctx, "user-123", "test-1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.Get(ctx, "user-123", "test-1")
	if err != ErrNotFound {
		t.Error("expected ErrNotFound after delete")
	}
}

func TestJSONStorageGetByUser(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{Type: "json", DataDir: tmpDir}
	store, _ := NewJSONStorage(cfg)
	defer store.Close()

	ctx := context.Background()

	// Сохраняем несколько записей
	entries := []*Entry{
		{ID: "1", UserID: "user-1", Section: "motivation", Answer: "a1"},
		{ID: "2", UserID: "user-1", Section: "boundaries", Answer: "a2"},
		{ID: "3", UserID: "user-2", Section: "motivation", Answer: "a3"},
	}

	for _, e := range entries {
		store.Save(ctx, e)
	}

	// Получаем записи user-1
	userEntries, _ := store.GetByUser(ctx, "user-1")
	if len(userEntries) != 2 {
		t.Errorf("expected 2 entries for user-1, got %d", len(userEntries))
	}
}

func TestJSONStorageGetBySection(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{Type: "json", DataDir: tmpDir}
	store, _ := NewJSONStorage(cfg)
	defer store.Close()

	ctx := context.Background()

	store.Save(ctx, &Entry{ID: "1", UserID: "user-1", Section: "motivation", Answer: "m1"})
	store.Save(ctx, &Entry{ID: "2", UserID: "user-1", Section: "boundaries", Answer: "b1"})
	store.Save(ctx, &Entry{ID: "3", UserID: "user-1", Section: "motivation", Answer: "m2"})

	entries, _ := store.GetBySection(ctx, "user-1", "motivation")
	if len(entries) != 2 {
		t.Errorf("expected 2 motivation entries, got %d", len(entries))
	}
}

func TestJSONStorageGetByTag(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{Type: "json", DataDir: tmpDir}
	store, _ := NewJSONStorage(cfg)
	defer store.Close()

	ctx := context.Background()

	store.Save(ctx, &Entry{ID: "1", UserID: "user-1", Tags: []string{"important", "work"}})
	store.Save(ctx, &Entry{ID: "2", UserID: "user-1", Tags: []string{"personal"}})
	store.Save(ctx, &Entry{ID: "3", UserID: "user-1", Tags: []string{"important", "health"}})

	entries, _ := store.GetByTag(ctx, "user-1", "important")
	if len(entries) != 2 {
		t.Errorf("expected 2 entries with tag important, got %d", len(entries))
	}
}

func TestJSONStorageExportImport(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{Type: "json", DataDir: tmpDir}
	store, _ := NewJSONStorage(cfg)
	defer store.Close()

	ctx := context.Background()

	// Экспортируем
	store.Save(ctx, &Entry{ID: "1", UserID: "user-1", Section: "motivation", Answer: "test"})
	entries, _ := store.Export(ctx, "user-1")

	if len(entries) != 1 {
		t.Errorf("expected 1 entry in export, got %d", len(entries))
	}

	// Импортируем
	err := store.Import(ctx, "user-2", entries)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	imported, _ := store.GetByUser(ctx, "user-2")
	if len(imported) != 1 {
		t.Errorf("expected 1 imported entry, got %d", len(imported))
	}
}

func TestJSONStorageCount(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{Type: "json", DataDir: tmpDir}
	store, _ := NewJSONStorage(cfg)
	defer store.Close()

	ctx := context.Background()

	count, _ := store.Count(ctx, "user-1")
	if count != 0 {
		t.Errorf("expected 0 entries initially, got %d", count)
	}

	store.Save(ctx, &Entry{ID: "1", UserID: "user-1", Section: "motivation"})
	store.Save(ctx, &Entry{ID: "2", UserID: "user-1", Section: "boundaries"})

	count, _ = store.Count(ctx, "user-1")
	if count != 2 {
		t.Errorf("expected 2 entries, got %d", count)
	}
}

func TestJSONStorageStats(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{Type: "json", DataDir: tmpDir}
	store, _ := NewJSONStorage(cfg)
	defer store.Close()

	ctx := context.Background()

	store.Save(ctx, &Entry{ID: "1", UserID: "user-1", Section: "motivation", Tags: []string{"tag1"}})
	store.Save(ctx, &Entry{ID: "2", UserID: "user-2", Section: "boundaries", Tags: []string{"tag2"}})

	stats, err := store.GetStats(ctx)
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}

	if stats.TotalUsers != 2 {
		t.Errorf("expected 2 users, got %d", stats.TotalUsers)
	}
	if stats.TotalEntries != 2 {
		t.Errorf("expected 2 entries, got %d", stats.TotalEntries)
	}
	if len(stats.SectionStats) != 2 {
		t.Errorf("expected 2 sections, got %d", len(stats.SectionStats))
	}
}

func TestJSONStorageReadOnly(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{Type: "json", DataDir: tmpDir, ReadOnly: true}
	store, _ := NewJSONStorage(cfg)
	defer store.Close()

	ctx := context.Background()

	err := store.Save(ctx, &Entry{ID: "1", UserID: "user-1"})
	if err != ErrReadOnly {
		t.Errorf("expected ErrReadOnly, got %v", err)
	}
}

func TestJSONStoragePersistence(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{Type: "json", DataDir: tmpDir}

	// Создаём и сохраняем
	store1, _ := NewJSONStorage(cfg)
	ctx := context.Background()
	store1.Save(ctx, &Entry{ID: "1", UserID: "user-1", Section: "motivation", Answer: "persistent"})
	store1.Close()

	// Создаём новое хранилище с тем же каталогом
	store2, _ := NewJSONStorage(cfg)
	defer store2.Close()

	// Проверяем, что данные загрузились
	entries, _ := store2.GetByUser(ctx, "user-1")
	if len(entries) != 1 {
		t.Error("expected data to persist")
	}
	if entries[0].Answer != "persistent" {
		t.Errorf("expected persistent answer, got %s", entries[0].Answer)
	}
}

func TestStorageError(t *testing.T) {
	err := &StorageError{Code: "test", Message: "test error"}

	if err.Error() != "test error" {
		t.Errorf("unexpected error message: %s", err.Error())
	}

	errWithWrap := &StorageError{Code: "test", Message: "test", Err: os.ErrNotExist}
	if errWithWrap.Unwrap() != os.ErrNotExist {
		t.Error("expected wrapped error")
	}
}

func TestRegistry(t *testing.T) {
	// Проверяем, что JSON зарегистрирован
	if _, ok := Registry["json"]; !ok {
		t.Error("expected json storage to be registered")
	}

	// Проверяем Factory
	factory := Registry["json"]
	if factory == nil {
		t.Error("expected non-nil factory")
	}
}

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{Type: "json", DataDir: tmpDir}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer store.Close()

	if store.Name() != "json" {
		t.Errorf("expected json storage, got %s", store.Name())
	}
}

func TestNewUnknownType(t *testing.T) {
	cfg := &Config{Type: "unknown", DataDir: "/tmp"}

	_, err := New(cfg)
	if err == nil {
		t.Error("expected error for unknown type")
	}
}

func TestConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{Type: "json", DataDir: tmpDir}
	store, _ := NewJSONStorage(cfg)
	defer store.Close()

	ctx := context.Background()
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			store.Save(ctx, &Entry{
				ID:       string(rune(id)),
				UserID:   "concurrent-user",
				Section:  "motivation",
				Answer:   "concurrent",
			})
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for concurrent operations")
		}
	}

	count, _ := store.Count(ctx, "concurrent-user")
	if count == 0 {
		t.Error("expected entries after concurrent access")
	}
}
