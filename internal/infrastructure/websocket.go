package infrastructure

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mu        sync.Mutex
	upgrader  websocket.Upgrader
}

var WebSocket = &WebSocketManager{
	clients:   make(map[*websocket.Conn]bool),
	broadcast: make(chan []byte),
	upgrader: websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	},
}

// Handle new WebSocket connections
func (wm *WebSocketManager) HandleConnections(c *gin.Context) {
	conn, err := wm.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	wm.mu.Lock()
	wm.clients[conn] = true
	wm.mu.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			wm.mu.Lock()
			delete(wm.clients, conn)
			wm.mu.Unlock()
			break
		}
		wm.broadcast <- msg
	}
}

// Start broadcasting messages to all clients
func (wm *WebSocketManager) StartBroadcaster() {
	for {
		msg := <-wm.broadcast
		wm.mu.Lock()
		for client := range wm.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				client.Close()
				delete(wm.clients, client)
			}
		}
		wm.mu.Unlock()
	}
}

// Send message to all connected clients
func (wm *WebSocketManager) SendMessage(pipelineName, stageName, status string) {
	message := map[string]string{
		"pipeline_name": pipelineName,
		"stage_name":    stageName,
		"status":        status,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling WebSocket message:", err)
		return
	}

	wm.broadcast <- jsonMessage
}
