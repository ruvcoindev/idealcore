// Package messages определяет типы сообщений для P2P-коммуникации
package messages

import (
	"fmt"
	"time"
)

// Message — базовая структура сообщения между нодами
type Message struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	From      string                 `json:"from"`
	To        string                 `json:"to,omitempty"`
	Timestamp int64                  `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
	Signature string                 `json:"signature,omitempty"`
	Priority  int                    `json:"priority,omitempty"`
	TTL       int                    `json:"ttl,omitempty"`
}

// Handshake — сообщение для установления соединения
type Handshake struct {
	NodeID    string                 `json:"node_id"`
	PublicKey string                 `json:"public_key"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// SyncRequest — запрос на синхронизацию данных
type SyncRequest struct {
	NodeID      string    `json:"node_id"`
	DataTypes   []string  `json:"data_types"`
	Since       time.Time `json:"since"`
	Compression bool      `json:"compression"`
}

// SyncResponse — ответ на запрос синхронизации
type SyncResponse struct {
	NodeID    string      `json:"node_id"`
	Data      interface{} `json:"data"`
	Count     int         `json:"count"`
	HasMore   bool        `json:"has_more"`
	Checksum  string      `json:"checksum"`
}

// DiaryEntrySync — синхронизация записи дневника
type DiaryEntrySync struct {
	UserID    string    `json:"user_id"`
	Section   string    `json:"section"`
	Answer    string    `json:"answer"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Hash      string    `json:"hash"`
}

// IntentionSync — синхронизация намерения
type IntentionSync struct {
	UserID      string    `json:"user_id"`
	Text        string    `json:"text"`
	Theme       string    `json:"theme"`
	CreatedAt   time.Time `json:"created_at"`
	UsedAt      time.Time `json:"used_at"`
	Hash        string    `json:"hash"`
}

// ConflictResolution — разрешение конфликтов при синхронизации
type ConflictResolution struct {
	Strategy    string    `json:"strategy"`
	LocalHash   string    `json:"local_hash"`
	RemoteHash  string    `json:"remote_hash"`
	ResolvedAt  time.Time `json:"resolved_at"`
}

// MessageType — типы сообщений
type MessageType string

const (
	MessageTypeHeartbeat      MessageType = "heartbeat"
	MessageTypeSyncRequest    MessageType = "sync_request"
	MessageTypeSyncResponse   MessageType = "sync_response"
	MessageTypeDiaryEntry     MessageType = "diary_entry"
	MessageTypeIntention      MessageType = "intention"
	MessageTypeConflict       MessageType = "conflict"
	MessageTypeAck            MessageType = "ack"
	MessageTypeError          MessageType = "error"
)

// NewMessage создаёт новое сообщение
func NewMessage(msgType MessageType, from string, payload map[string]interface{}) *Message {
	return &Message{
		ID:        generateMessageID(),
		Type:      string(msgType),
		From:      from,
		Timestamp: time.Now().Unix(),
		Payload:   payload,
		TTL:       5,
		Priority:  0,
	}
}

// generateMessageID генерирует уникальный ID сообщения
func generateMessageID() string {
	return fmt.Sprintf("msg-%d-%x", time.Now().UnixNano(), time.Now().UnixNano()%0xFFFFFF)
}

// Validate проверяет валидность сообщения
func (m *Message) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("message ID is required")
	}
	if m.Type == "" {
		return fmt.Errorf("message type is required")
	}
	if m.From == "" {
		return fmt.Errorf("message sender is required")
	}
	if m.Timestamp == 0 {
		return fmt.Errorf("message timestamp is required")
	}
	return nil
}

// IsExpired проверяет, не истёк ли TTL сообщения
func (m *Message) IsExpired() bool {
	if m.TTL <= 0 {
		return false
	}
	age := time.Now().Unix() - m.Timestamp
	return age > int64(m.TTL*60)
}

// DecrementTTL уменьшает счётчик TTL
func (m *Message) DecrementTTL() {
	if m.TTL > 0 {
		m.TTL--
	}
}
