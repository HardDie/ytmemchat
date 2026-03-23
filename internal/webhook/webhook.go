// Package webhook provides HTTP handlers to manually inject messages into the
// application pipeline or trigger system commands (like interrupting TTS).
// This is primarily used for testing, debugging, or integration with
// external third-party tools.
package webhook

import (
	"time"

	clientYoutube "github.com/HardDie/ytmemchat/internal/clients/youtube"
	"github.com/HardDie/ytmemchat/internal/server"

	"encoding/json"
	"net/http"
)

// Webhook manages the ingestion of external HTTP signals.
type Webhook struct {
	ch        chan *clientYoutube.ChatMessage
	broadcast chan server.WebsocketPayload
}

// New initializes a Webhook instance with the required broadcast channel.
func New(cfg Config) *Webhook {
	return &Webhook{
		ch:        make(chan *clientYoutube.ChatMessage),
		broadcast: cfg.Broadcast,
	}
}

// Handle accepts a JSON payload containing a "message" and injects it into
// the system as a mock YouTube ChatMessage.
//
// Payload format: {"message": "@wow"}
func (wh *Webhook) Handle(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Message string `json:"message"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	wh.ch <- &clientYoutube.ChatMessage{
		Timestamp: time.Now(),
		Type:      "debug",
		Author:    "webhook",
		Message:   payload.Message,
	}
}

// InterruptHandle immediately broadcasts a TTS interrupt signal to all
// connected WebSocket clients, stopping any currently playing speech.
func (wh *Webhook) InterruptHandle(_ http.ResponseWriter, _ *http.Request) {
	wh.broadcast <- server.WebsocketPayload{
		Type: server.PayloadTypeTTSInterrupt,
	}
}

// GetChan returns the channel where injected messages are published.
func (wh *Webhook) GetChan() chan *clientYoutube.ChatMessage {
	return wh.ch
}
