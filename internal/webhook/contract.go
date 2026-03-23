package webhook

import "github.com/HardDie/ytmemchat/internal/server"

type Config struct {
	Broadcast chan server.WebsocketPayload
}
