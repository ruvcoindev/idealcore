// Package yggdrasil — P2P транспорт (заглушка для сборки)
package yggdrasil

// Transport — P2P транспорт
type Transport struct {
	nodeID  string
	running bool
}

// NewTransport создаёт транспорт
func NewTransport(cfg interface{}) *Transport {
	return &Transport{
		nodeID:  "idealcore-local",
		running: false,
	}
}

// Start запускает транспорт
func (t *Transport) Start() error {
	t.running = true
	return nil
}

// Stop останавливает транспорт
func (t *Transport) Stop() error {
	t.running = false
	return nil
}

// GetNodeID возвращает ID ноды
func (t *Transport) GetNodeID() string {
	return t.nodeID
}

// IsRunning возвращает статус
func (t *Transport) IsRunning() bool {
	return t.running
}

// Peer — информация о пире
type Peer struct {
	ID string
}

// GetPeers возвращает список пиров
func (t *Transport) GetPeers() []*Peer {
	return nil
}

// TransportStats — статистика
type TransportStats struct {
	NodeID    string `json:"node_id"`
	PeerCount int    `json:"peer_count"`
	Running   bool   `json:"running"`
}

// GetStats возвращает статистику
func (t *Transport) GetStats() TransportStats {
	return TransportStats{
		NodeID:    t.nodeID,
		PeerCount: 0,
		Running:   t.running,
	}
}
