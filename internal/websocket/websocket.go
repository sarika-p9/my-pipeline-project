package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Upgrader settings for WebSocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// StageUpdate represents the message structure for stage updates
type StageUpdate struct {
	StageID string `json:"stage_id"`
	Status  string `json:"status"`
}

// WebSocketManager handles WebSocket clients
type WebSocketManager struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

// Manager is a global WebSocket manager instance
var Manager = &WebSocketManager{
	clients: make(map[*websocket.Conn]bool),
}

// HandleConnections upgrades HTTP connections to WebSockets
func (wm *WebSocketManager) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[WebSocket] Upgrade failed:", err)
		http.Error(w, "WebSocket Upgrade failed", http.StatusInternalServerError)
		return
	}

	// Add client to map
	wm.mu.Lock()
	wm.clients[conn] = true
	wm.mu.Unlock()

	log.Println("[WebSocket] New client connected")

	// Ensure cleanup on disconnect
	defer wm.removeClient(conn)

	// Start ping mechanism in a separate goroutine
	go wm.keepAlive(conn)

	// Keep reading messages to detect disconnects
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("[WebSocket] Read error, disconnecting client:", err)
			break
		}
	}
}

// keepAlive sends periodic pings to maintain the WebSocket connection
func (wm *WebSocketManager) keepAlive(conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Println("[WebSocket] Ping failed, disconnecting client:", err)
			wm.removeClient(conn)
			return
		}
	}
}

// removeClient safely removes a client from the connection map
func (wm *WebSocketManager) removeClient(conn *websocket.Conn) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	delete(wm.clients, conn)
	conn.Close()
	log.Println("[WebSocket] Client disconnected and removed")
}

// BroadcastMessage sends JSON messages to all connected clients
func (wm *WebSocketManager) BroadcastMessage(stageID, status string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	message := StageUpdate{
		StageID: stageID,
		Status:  status,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println("[WebSocket] Failed to marshal JSON:", err)
		return
	}

	for client := range wm.clients {
		if err := client.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
			log.Println("[WebSocket] Write failed, removing client:", err)
			wm.removeClient(client)
		}
	}
}
