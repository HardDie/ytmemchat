package server

import (
	"log/slog"
	"net/http"

	"github.com/HardDie/ytmemchat/pkg/logger"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader settings. CheckOrigin allows the connection from OBS.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity in this local OBS setup
	},
}

var wsClients = make(map[*websocket.Conn]bool) // Connected clients (WebSocket connections)

// wsHandler manages new WebSocket connections.
func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error(
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
			s.logger.Error(
				"conn.ReadMessage()",
				slog.String(logger.LogValueError, err.Error()),
			)
			delete(wsClients, conn) // Unregister client on error/close
			return
		}
	}
}
