package server

type Config struct {
	Port string
}

type PayloadType string

const (
	PayloadTypeAlert PayloadType = "alert"
	PayloadTypeTTS   PayloadType = "tts"
)

// WebsocketPayload defines the structure of the websocket payload JSON.
type WebsocketPayload struct {
	Type     PayloadType `json:"type"`
	Filename string      `json:"filename"` // e.g., "new_follower.gif", "cheer_alert.mp4"
	Volume   float64     `json:"volume"`   // Field for volume (0.0 to 1.0)
	Scale    float64     `json:"scale"`    // Field for scale multiplier (e.g., 0.5 to 2.0)
}
