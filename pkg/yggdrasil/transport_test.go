package yggdrasil

import (
	"testing"
	"time"

	"github.com/ruvcoindev/idealcore/pkg/config"
	"github.com/ruvcoindev/idealcore/pkg/yggdrasil/messages"
)

func TestNewTransport(t *testing.T) {
	cfg := config.Load()
	transport := NewTransport(cfg)

	if transport == nil {
		t.Fatal("expected transport, got nil")
	}
	if transport.nodeID == "" {
		t.Error("expected nodeID to be set")
	}
	if transport.running {
		t.Error("expected transport to be stopped initially")
	}
}

func TestTransportStartStop(t *testing.T) {
	cfg := config.Load()
	transport := NewTransport(cfg)

	// Запускаем
	err := transport.Start()
	if err != nil {
		t.Fatalf("failed to start transport: %v", err)
	}

	// Даём время на запуск
	time.Sleep(100 * time.Millisecond)

	if !transport.IsRunning() {
		t.Error("expected transport to be running")
	}

	// Останавливаем
	err = transport.Stop()
	if err != nil {
		t.Fatalf("failed to stop transport: %v", err)
	}

	if transport.IsRunning() {
		t.Error("expected transport to be stopped")
	}
}

func TestGetNodeID(t *testing.T) {
	cfg := config.Load()
	transport := NewTransport(cfg)

	nodeID := transport.GetNodeID()
	if nodeID == "" {
		t.Error("expected non-empty nodeID")
	}

	// Проверяем формат
	expectedPrefix := "idealcore-"
	if nodeID[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("expected nodeID to start with %s, got %s", expectedPrefix, nodeID)
	}
}

func TestGetPeers(t *testing.T) {
	cfg := config.Load()
	transport := NewTransport(cfg)

	peers := transport.GetPeers()
	if len(peers) != 0 {
		t.Errorf("expected 0 peers initially, got %d", len(peers))
	}
}

func TestGetStats(t *testing.T) {
	cfg := config.Load()
	transport := NewTransport(cfg)

	stats := transport.GetStats()

	if stats.NodeID == "" {
		t.Error("expected non-empty nodeID in stats")
	}
	if stats.PeerCount != 0 {
		t.Errorf("expected 0 peers in stats, got %d", stats.PeerCount)
	}
	if stats.Running {
		t.Error("expected running to be false initially")
	}
}

func TestSubscribe(t *testing.T) {
	cfg := config.Load()
	transport := NewTransport(cfg)

	called := false
	handler := func(msg *messages.Message) error {
		called = true
		return nil
	}

	transport.Subscribe("test_message", handler)
	transport.Subscribe("test_message", handler) // подписываемся дважды

	stats := transport.GetStats()
	if stats.SubscriberCount < 1 {
		t.Errorf("expected at least 1 subscriber, got %d", stats.SubscriberCount)
	}
}

func TestMessageValidation(t *testing.T) {
	msg := &messages.Message{
		ID:        "test-1",
		Type:      "test",
		From:      "node-1",
		Timestamp: time.Now().Unix(),
		Payload:   map[string]interface{}{"key": "value"},
	}

	err := msg.Validate()
	if err != nil {
		t.Errorf("expected valid message, got error: %v", err)
	}
}

func TestMessageExpiration(t *testing.T) {
	msg := &messages.Message{
		ID:        "test-2",
		Type:      "test",
		From:      "node-1",
		Timestamp: time.Now().Unix() - 1000, // 1000 секунд назад
		TTL:       5, // 5 минут
		Payload:   map[string]interface{}{},
	}

	if !msg.IsExpired() {
		t.Error("expected message to be expired")
	}
}

func TestMessageTTL(t *testing.T) {
	msg := &messages.Message{
		ID:        "test-3",
		Type:      "test",
		From:      "node-1",
		Timestamp: time.Now().Unix(),
		TTL:       5,
		Payload:   map[string]interface{}{},
	}

	initialTTL := msg.TTL
	msg.DecrementTTL()

	if msg.TTL != initialTTL-1 {
		t.Errorf("expected TTL to decrease from %d to %d, got %d", initialTTL, initialTTL-1, msg.TTL)
	}
}

func TestBroadcast(t *testing.T) {
	cfg := config.Load()
	transport := NewTransport(cfg)

	msg := messages.NewMessage(
		messages.MessageTypeHeartbeat,
		transport.GetNodeID(),
		map[string]interface{}{"status": "alive"},
	)

	// Broadcast без пиров не должен паниковать
	err := transport.Broadcast(msg)
	if err != nil {
		t.Logf("broadcast error (expected with no peers): %v", err)
	}
}
