package server

// Config defines the network settings for the server.
type Config struct {
	Port string // The address/port to bind to (e.g., ":8080")
}

type PayloadType string

const (
	PayloadTypeAlert        PayloadType = "alert"
	PayloadTypeTTS          PayloadType = "tts"
	PayloadTypeTTSInterrupt PayloadType = "tts_interrupt"
)

// WebsocketPayload defines the JSON structure sent to the browser.
type WebsocketPayload struct {
	Type     PayloadType `json:"type"`     // "alert", "tts", or "tts_interrupt"
	Payload  []byte      `json:"payload"`  // Raw audio data (used for TTS)
	Filename string      `json:"filename"` // Filename in the media directory (for alerts)
	Volume   float64     `json:"volume"`   // 0.0 to 1.0
	Scale    float64     `json:"scale"`    // Visual size multiplier
}
