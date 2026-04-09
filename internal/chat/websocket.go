package chat

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/HardDie/ytmemchat/pkg/logger"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var wsClients = make(map[*websocket.Conn]bool) // Connected clients (WebSocket connections)

// WSHandler manages new WebSocket connections.
func (c *Chat) WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.logger.Error(
			"upgrader.Upgrade()",
			slog.String(logger.LogValueError, err.Error()),
		)
		return
	}
	defer conn.Close()

	wsClients[conn] = true // Register new client

	// Keep the connection alive by reading (though we expect no messages from client)
	for {
		_, _, err = conn.ReadMessage()
		if err != nil {
			if !strings.Contains(err.Error(), "going away") {
				c.logger.Error(
					"conn.ReadMessage()",
					slog.String(logger.LogValueError, err.Error()),
				)
			}
			delete(wsClients, conn) // Unregister client on error/close
			return
		}
	}
}
