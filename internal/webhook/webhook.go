package webhook

import (
	"time"

	clientYoutube "github.com/HardDie/ytmemchat/internal/clients/youtube"
	"github.com/HardDie/ytmemchat/internal/server"

	"encoding/json"
	"net/http"
)

type Webhook struct {
	ch        chan *clientYoutube.ChatMessage
	broadcast chan server.WebsocketPayload
}

func New(cfg Config) *Webhook {
	return &Webhook{
		ch:        make(chan *clientYoutube.ChatMessage),
		broadcast: cfg.Broadcast,
	}
}

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

func (wh *Webhook) InterruptHandle(_ http.ResponseWriter, _ *http.Request) {
	wh.broadcast <- server.WebsocketPayload{
		Type: server.PayloadTypeTTSInterrupt,
	}
}

func (wh *Webhook) GetChan() chan *clientYoutube.ChatMessage {
	return wh.ch
}
