package diary

import (
	"os"
	"testing"
	"time"
)

func TestNewStore(t *testing.T) {
	tmpDir := t.TempDir()
	
	store, err := NewStore(tmpDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	if store == nil {
		t.Fatal("expected store, got nil")
	}
	if store.dataDir != tmpDir {
		t.Errorf("expected dataDir %s, got %s", tmpDir, store.dataDir)
	}
}

func TestSaveEntry(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewStore(tmpDir)

	userID := "test-user"
	section := "motivation"
	answer := "Я делаю это для себя, а не для других"
	tags := []string{"motivation", "self"}

	err := store.SaveEntry(userID, section, answer, tags)
	if err != nil {
		t.Fatalf("failed to save entry: %v", err)
	}

	// Проверяем, что запись сохранилась
	entries, err := store.GetEntries(userID)
	if err != nil {
		t.Fatalf("failed to get entries: %v", err)
	}

	if len(entries) == 0 {
		t.Error("expected at least one entry")
	}
}

func TestSaveEntryValidation(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewStore(tmpDir)

	// Пустой userID
	err := store.SaveEntry("", "motivation", "answer", nil)
	if err == nil {
		t.Error("expected error for empty userID")
	}

	// Пустая секция
	err = store.SaveEntry("user", "", "answer", nil)
	if err == nil {
		t.Error("expected error for empty section")
	}
}

func TestGetEntries(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewStore(tmpDir)

	userID := "test-user-2"
	
	// Сохраняем запись
	_ = store.SaveEntry(userID, "motivation", "test answer", []string{"test"})

	// Получаем записи
	entries, err := store.GetEntries(userID)
	if err != nil {
		t.Fatalf("failed to get entries: %v", err)
	}

	if len(entries) == 0 {
		t.Error("expected at least one entry")
	}
}

func TestGetEntry(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewStore(tmpDir)

	userID := "test-user-3"
	_ = store.SaveEntry(userID, "motivation", "test answer", []string{"test"})

	entry, ok := store.GetEntry(userID, "motivation_1")
	if !ok {
		t.Fatal("expected to find entry")
	}

	if entry.Answer != "test answer" {
		t.Errorf("expected answer 'test answer', got %s", entry.Answer)
	}
}

func TestDeleteEntry(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewStore(tmpDir)

	userID := "test-user-4"
	_ = store.SaveEntry(userID, "motivation", "test answer", []string{"test"})

	err := store.DeleteEntry(userID, "motivation_1")
	if err != nil {
		t.Fatalf("failed to delete entry: %v", err)
	}

	// Проверяем, что удалено
	entries, _ := store.GetEntries(userID)
	if len(entries) > 0 {
		t.Error("expected entry to be deleted")
	}
}

func TestGetEntriesBySection(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewStore(tmpDir)

	userID := "test-user-5"
	_ = store.SaveEntry(userID, "motivation", "motivation answer", []string{"motivation"})
	_ = store.SaveEntry(userID, "boundaries", "boundaries answer", []string{"boundaries"})

	entries, err := store.GetEntriesBySection(userID, "motivation")
	if err != nil {
		t.Fatalf("failed to get entries by section: %v", err)
	}

	if len(entries) == 0 {
		t.Error("expected at least one motivation entry")
	}

	for _, entry := range entries {
		if entry.Section != "motivation" {
			t.Errorf("expected section motivation, got %s", entry.Section)
		}
	}
}

func TestGetEntriesByTag(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewStore(tmpDir)

	userID := "test-user-6"
	_ = store.SaveEntry(userID, "motivation", "answer 1", []string{"important", "motivation"})
	_ = store.SaveEntry(userID, "boundaries", "answer 2", []string{"boundaries"})

	entries, err := store.GetEntriesByTag(userID, "important")
	if err != nil {
		t.Fatalf("failed to get entries by tag: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("expected 1 entry with tag 'important', got %d", len(entries))
	}
}

func TestStorePersistence(t *testing.T) {
	tmpDir := t.TempDir()

	// Создаём хранилище и сохраняем данные
	store1, _ := NewStore(tmpDir)
	userID := "persist-user"
	_ = store1.SaveEntry(userID, "motivation", "persistent answer", []string{"test"})

	// Создаём новое хранилище с тем же каталогом
	store2, _ := NewStore(tmpDir)

	// Проверяем, что данные загрузились
	entries, err := store2.GetEntries(userID)
	if err != nil {
		t.Fatalf("failed to get entries: %v", err)
	}

	if len(entries) == 0 {
		t.Error("expected data to persist across store instances")
	}
}

func TestCount(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewStore(tmpDir)

	userID := "count-user"
	
	// Изначально 0
	if store.Count(userID) != 0 {
		t.Error("expected 0 entries initially")
	}

	// Сохраняем запись
	_ = store.SaveEntry(userID, "motivation", "answer", []string{"test"})

	// Проверяем счётчик
	count := store.Count(userID)
	if count == 0 {
		t.Error("expected count > 0 after save")
	}
}

func TestGetStats(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewStore(tmpDir)

	userID := "stats-user"
	_ = store.SaveEntry(userID, "motivation", "answer 1", []string{"test"})
	_ = store.SaveEntry(userID, "boundaries", "answer 2", []string{"test"})

	stats := store.GetStats()

	if stats.TotalUsers < 1 {
		t.Error("expected at least 1 user")
	}
	if stats.TotalEntries < 2 {
		t.Error("expected at least 2 entries")
	}
	if len(stats.SectionStats) < 2 {
		t.Error("expected at least 2 sections in stats")
	}
}

func TestExportImport(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewStore(tmpDir)

	userID := "export-user"
	_ = store.SaveEntry(userID, "motivation", "export answer", []string{"test"})

	// Экспортируем
	data, err := store.Export(userID)
	if err != nil {
		t.Fatalf("failed to export: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty export data")
	}

	// Импортируем в нового пользователя
	importUserID := "import-user"
	err = store.Import(importUserID, data)
	if err != nil {
		t.Fatalf("failed to import: %v", err)
	}

	// Проверяем, что данные импортировались
	importedEntries, _ := store.GetEntries(importUserID)
	if len(importedEntries) == 0 {
		t.Error("expected imported entries")
	}
}

func TestConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	store, _ := NewStore(tmpDir)

	userID := "concurrent-user"
	done := make(chan bool, 10)

	// Запускаем 10 горутин для записи
	for i := 0; i < 10; i++ {
		go func(id int) {
			_ = store.SaveEntry(userID, "motivation", "concurrent answer", []string{"test"})
			done <- true
		}(i)
	}

	// Ждём завершения
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for concurrent operations")
		}
	}

	// Проверяем, что данные сохранились
	entries, _ := store.GetEntries(userID)
	if len(entries) == 0 {
		t.Error("expected entries after concurrent access")
	}
}
