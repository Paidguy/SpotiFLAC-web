package server

import (
	"encoding/json"
	"fmt"
	"sync"
)

// SSEClient represents a connected client
type SSEClient struct {
	ID      string
	Channel chan []byte
}

// SSEBroker manages SSE connections and broadcasts messages
type SSEBroker struct {
	clients    map[string]*SSEClient
	clientsMux sync.RWMutex
	register   chan *SSEClient
	unregister chan *SSEClient
	broadcast  chan []byte
}

// NewSSEBroker creates a new SSE broker
func NewSSEBroker() *SSEBroker {
	return &SSEBroker{
		clients:    make(map[string]*SSEClient),
		register:   make(chan *SSEClient),
		unregister: make(chan *SSEClient),
		broadcast:  make(chan []byte, 256),
	}
}

// Run starts the broker's event loop
func (b *SSEBroker) Run() {
	for {
		select {
		case client := <-b.register:
			b.clientsMux.Lock()
			b.clients[client.ID] = client
			b.clientsMux.Unlock()

		case client := <-b.unregister:
			b.clientsMux.Lock()
			if _, ok := b.clients[client.ID]; ok {
				delete(b.clients, client.ID)
				close(client.Channel)
			}
			b.clientsMux.Unlock()

		case message := <-b.broadcast:
			b.clientsMux.RLock()
			for _, client := range b.clients {
				select {
				case client.Channel <- message:
				default:
					// Client is slow or disconnected, skip
				}
			}
			b.clientsMux.RUnlock()
		}
	}
}

// RegisterClient registers a new SSE client
func (b *SSEBroker) RegisterClient(client *SSEClient) {
	b.register <- client
}

// UnregisterClient removes a client
func (b *SSEBroker) UnregisterClient(client *SSEClient) {
	b.unregister <- client
}

// Broadcast sends a message to all connected clients
func (b *SSEBroker) Broadcast(message []byte) {
	b.broadcast <- message
}

// BroadcastJSON sends a JSON message to all connected clients
func (b *SSEBroker) BroadcastJSON(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	b.Broadcast(jsonData)
	return nil
}
