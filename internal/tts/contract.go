package tts

import "github.com/HardDie/ytmemchat/internal/server"

type Config struct {
	VoiceName string
	Broadcast chan server.WebsocketPayload
	Volume    *float64
}
