package youtube

import (
	"context"
	"time"
)

// ChatMessage represents the structured data returned by the iterator.
// This is independent of the YouTube API structure.
type ChatMessage struct {
	ID        string
	Author    string
	Message   string
	Type      string // e.g., "textMessageEvent", "superChatEvent"
	Timestamp time.Time
}

// Client defines the interface for any live chat source.
// This is the contract your core application depends on.
type Client interface {
	// GetMessageIterator returns an iterator that will continuously poll for new messages.
	GetMessageIterator(ctx context.Context, liveVideoID string) (MessageIterator, error)
}

// MessageIterator is the core abstraction for consuming messages.
type MessageIterator interface {
	// Next returns the next message. It blocks until a message is available.
	// It returns (nil, false) if the context is cancelled or the stream ends.
	Next() (*ChatMessage, bool)
	GetChan() chan *ChatMessage
}
