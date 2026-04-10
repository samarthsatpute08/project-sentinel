package telemetry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type Metrics struct {
	Timestamp      int64   `json:"timestamp"`
	RPS            float64 `json:"rps"`
	ActiveRoute    string  `json:"activeRoute"` // "primary" or "fallback"
	CircuitState   string  `json:"circuitState"`
	TotalRequests  int64   `json:"totalRequests"`
	TotalFailures  int64   `json:"totalFailures"`
	Latency        float64 `json:"latency"`
}

type Hub struct {
	clients   map[*websocket.Conn]struct{}
	mu        sync.RWMutex
	broadcast chan []byte
}

func NewHub() *Hub {
	h := &Hub{
		clients:   make(map[*websocket.Conn]struct{}),
		broadcast: make(chan []byte, 256), // buffered - never blocks the main path
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for msg := range h.broadcast {
		h.mu.RLock()
		for conn := range h.clients {
			// Non-blocking write with deadline to avoid slow client blocking hub
			conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
			if err := websocket.Message.Send(conn, string(msg)); err != nil {
				// Will be cleaned up on next read error
				_ = err
			}
		}
		h.mu.RUnlock()
	}
}

func (h *Hub) Publish(m Metrics) {
	m.Timestamp = time.Now().UnixMilli()
	data, err := json.Marshal(m)
	if err != nil {
		return
	}
	// Non-blocking send - if buffer is full, drop the metric (it's telemetry, not critical)
	select {
	case h.broadcast <- data:
	default:
		log.Println("telemetry: broadcast buffer full, dropping metric")
	}
}

func (h *Hub) Handler() http.Handler {
	return websocket.Handler(func(conn *websocket.Conn) {
		h.mu.Lock()
		h.clients[conn] = struct{}{}
		h.mu.Unlock()

		log.Printf("telemetry: client connected (%d total)", len(h.clients))

		// Block until client disconnects (read loop)
		buf := make([]byte, 64)
		for {
			_, err := conn.Read(buf)
			if err != nil {
				break
			}
		}

		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		log.Printf("telemetry: client disconnected (%d total)", len(h.clients))
	})
}