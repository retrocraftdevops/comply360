package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Configure appropriately for production
	},
}

// Client represents a WebSocket client
type Client struct {
	ID       string
	TenantID string
	Conn     *websocket.Conn
	Send     chan []byte
}

// Hub maintains active WebSocket connections
type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

// Global hub instance
var hub = &Hub{
	Clients:    make(map[string]*Client),
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

// InitializeHub starts the WebSocket hub
func InitializeHub() {
	go hub.Run()
	go hub.BroadcastMetrics()
}

// Run handles WebSocket hub operations
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
			log.Printf("Client registered: %s (Tenant: %s)", client.ID, client.TenantID)
			
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: %s", client.ID)
			
		case message := <-h.Broadcast:
			h.mu.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.ID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastMetrics sends periodic metric updates to all connected clients
func (h *Hub) BroadcastMetrics() {
	ticker := time.NewTicker(5 * time.Second) // Update every 5 seconds
	defer ticker.Stop()
	
	for {
		<-ticker.C
		
		h.mu.RLock()
		if len(h.Clients) == 0 {
			h.mu.RUnlock()
			continue
		}
		h.mu.RUnlock()
		
		// Generate metric update
		metrics := generateMetricUpdate()
		
		// Marshal to JSON
		data, err := json.Marshal(metrics)
		if err != nil {
			log.Printf("Error marshaling metrics: %v", err)
			continue
		}
		
		// Broadcast to all clients
		h.Broadcast <- data
	}
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-ID")
	// System users (system_admin, global_admin) may have empty tenant_id

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	
	// Create client
	client := &Client{
		ID:       generateClientID(),
		TenantID: tenantID,
		Conn:     conn,
		Send:     make(chan []byte, 256),
	}
	
	// Register client
	hub.Register <- client
	
	// Start goroutines for reading and writing
	go client.ReadPump()
	go client.WritePump()
}

// ReadPump reads messages from the WebSocket
func (c *Client) ReadPump() {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		
		// Handle incoming messages (subscriptions, etc.)
		log.Printf("Received message from client %s: %s", c.ID, message)
	}
}

// WritePump writes messages to the WebSocket
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			
			// Add queued messages to the current WebSocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}
			
			if err := w.Close(); err != nil {
				return
			}
			
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Helper functions
func generateClientID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

func generateMetricUpdate() map[string]interface{} {
	return map[string]interface{}{
		"type":      "metric_update",
		"timestamp": time.Now().Unix(),
		"data": map[string]interface{}{
			"total_mrr":           125000 + (time.Now().Unix() % 1000),
			"active_users_now":    45 + (time.Now().Unix() % 10),
			"registrations_today": 3,
			"logins_today":        127 + (time.Now().Unix() % 20),
		},
	}
}

