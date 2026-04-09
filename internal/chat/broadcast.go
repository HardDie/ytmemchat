package chat

import (
	"log/slog"

	"github.com/HardDie/ytmemchat/pkg/logger"
)

var broadcast = make(chan WebsocketPayload) // Channel to broadcast webhook data to clients

// broadcaster constantly monitors the 'broadcast' channel and sends messages to all clients.
func (c *Chat) broadcaster() {
	for {
		message := <-broadcast // Wait for a message
		// Send the received payload to every connected client
		for client := range wsClients {
			err := client.WriteJSON(message)
			if err != nil {
				c.logger.Error(
					"client.WriteJSON()",
					slog.String(logger.LogValueError, err.Error()),
				)
				client.Close()
				delete(wsClients, client)
			}
		}
	}
}
