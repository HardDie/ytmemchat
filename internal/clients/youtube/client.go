// Package youtube provides a client for interacting with YouTube Live Chat.
// It abstracts the polling logic, pagination, and API-specific data structures
// into a simple, channel-based iterator.
package youtube

import (
	"context"
	"time"
)

// ChatMessage represents a normalized chat event.
// It abstracts away the differences between regular text messages,
// Super Chats, and other event types.
type ChatMessage struct {
	ID        string
	Author    string
	ImgURL    string
	Message   string // The text content or a formatted string for Super Chats
	Type      string // The original YouTube event type (e.g., "textMessageEvent")
	Timestamp time.Time
}

// Client is the primary interface for the YouTube chat service.
type Client interface {
	// GetMessageIterator initializes a chat session for a specific video
	// and returns an iterator to consume messages.
	GetMessageIterator(ctx context.Context, liveVideoID string) (MessageIterator, error)
}

// MessageIterator provides a way to consume chat messages sequentially.
type MessageIterator interface {
	// Next blocks until a new message is available or the context is closed.
	// Returns (nil, false) when the stream ends or is cancelled.
	Next() (*ChatMessage, bool)
	// GetChan returns the underlying channel for use in select statements.
	GetChan() chan *ChatMessage
}
