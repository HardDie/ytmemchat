package chat

import (
	"log/slog"

	clientYoutube "github.com/HardDie/ytmemchat/internal/clients/youtube"
	"github.com/HardDie/ytmemchat/pkg/logger"
)

type Chat struct {
	logger *slog.Logger
}

func New() *Chat {

	c := &Chat{
		logger: logger.Logger.With(slog.String(logger.LogService, "chat")),
	}

	// Start the broadcaster in a goroutine
	go c.broadcaster()

	return c
}

func (c *Chat) Message(msg *clientYoutube.ChatMessage) {
	// html.EscapeString()
	broadcast <- WebsocketPayload{
		AuthorName:    msg.Author,
		AuthorPicture: msg.ImgURL,
		MessageText:   msg.Message,
		PublishedAt:   msg.Timestamp.String(),
		IsModerator:   false,
		IsOwner:       false,
	}
}
