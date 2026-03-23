package tts

import "github.com/HardDie/ytmemchat/internal/server"

// Config defines the settings for the TTS engine.
type Config struct {
	VoiceName string                       // The system name of the voice to use
	Broadcast chan server.WebsocketPayload // Channel to send audio data to the frontend
	Volume    *float64                     // Optional: Playback volume (0.0 to 1.0)
}
