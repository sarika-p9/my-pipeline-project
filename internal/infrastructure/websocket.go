package infrastructure

import (
	"fmt"
	"sync"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketManager manages all active connections
type WebSocketManager struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan string
	Mutex     sync.Mutex
}

var WebSocket = &WebSocketManager{
	Clients:   make(map[*websocket.Conn]bool),
	Broadcast: make(chan string),
}

// WebSocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocket handler function
func (wm *WebSocketManager) HandleConnections(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("❌ Failed to upgrade WebSocket connection:", err)
		return
	}
	defer conn.Close()

	wm.Mutex.Lock()
	wm.Clients[conn] = true
	wm.Mutex.Unlock()

	fmt.Println("✅ New WebSocket client connected")

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			wm.Mutex.Lock()
			delete(wm.Clients, conn)
			wm.Mutex.Unlock()
			fmt.Println("❌ WebSocket client disconnected")
			break
		}
	}
}

// Broadcast messages to all clients
func (wm *WebSocketManager) BroadcastMessage(message string) {
	wm.Mutex.Lock()
	defer wm.Mutex.Unlock()

	for client := range wm.Clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("❌ Failed to send message to WebSocket client:", err)
			client.Close()
			delete(wm.Clients, client)
		}
	}
}

// Start the WebSocket broadcast listener
func (wm *WebSocketManager) StartBroadcaster() {
	go func() {
		for {
			msg := <-wm.Broadcast
			wm.BroadcastMessage(msg)
		}
	}()
}
