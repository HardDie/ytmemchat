package alerts

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader settings. CheckOrigin allows the connection from OBS.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity in this local OBS setup
	},
}

var clients = make(map[*websocket.Conn]bool) // Connected clients (WebSocket connections)
var broadcast = make(chan WebhookPayload)    // Channel to broadcast webhook data to clients

// wsHandler manages new WebSocket connections.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true // Register new client

	// Keep the connection alive by reading (though we expect no messages from client)
	for {
		_, _, err = conn.ReadMessage()
		if err != nil {
			log.Printf("Read error for client %v: %v", conn.RemoteAddr(), err)
			delete(clients, conn) // Unregister client on error/close
			return
		}
	}
}

// broadcaster constantly monitors the 'broadcast' channel and sends messages to all clients.
func broadcaster() {
	for {
		message := <-broadcast // Wait for a message
		// Send the received payload to every connected client
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("Write error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
