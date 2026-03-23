package webhook

import (
	"time"

	clientYoutube "github.com/HardDie/ytmemchat/internal/clients/youtube"

	"encoding/json"
	"net/http"
)

type Webhook struct {
	ch chan *clientYoutube.ChatMessage
}

func New() *Webhook {
	return &Webhook{
		ch: make(chan *clientYoutube.ChatMessage),
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

func (wh *Webhook) GetChan() chan *clientYoutube.ChatMessage {
	return wh.ch
}
