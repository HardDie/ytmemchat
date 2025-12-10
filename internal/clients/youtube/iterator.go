package youtube

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/HardDie/ytmemchat/pkg/logger"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// --- Iterator Implementation ---

type youtubeIterator struct {
	ctx          context.Context
	service      *youtube.Service
	liveChatID   string
	pageToken    string
	pollingDelay time.Duration
	messageChan  chan *ChatMessage
	logger       *slog.Logger
}

// New creates and initializes the YouTube API client.
func New(apiKey string) (Client, error) {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube service: %w", err)
	}
	return &youtubeAPIClient{
		service: service,
		logger:  logger.Logger.With(slog.String(logger.LogService, "YouTube Client")),
	}, nil
}

type youtubeAPIClient struct {
	service *youtube.Service
	logger  *slog.Logger
}

func (c *youtubeAPIClient) GetMessageIterator(ctx context.Context, liveVideoID string) (MessageIterator, error) {
	liveChatID, err := c.getLiveChatID(liveVideoID)
	if err != nil {
		return nil, err
	}

	it := &youtubeIterator{
		ctx:          ctx,
		service:      c.service,
		liveChatID:   liveChatID,
		pageToken:    "",              // Start with an empty token
		pollingDelay: 5 * time.Second, // Default initial delay
		messageChan:  make(chan *ChatMessage, 10),
		logger:       c.logger,
	}

	// Perform the initial setup poll to discard old messages and get the starting token.
	if err = it.initializeToken(); err != nil {
		return nil, fmt.Errorf("failed to initialize chat token: %w", err)
	}

	// Start the background goroutine to handle polling and channel population
	go it.startPolling()

	return it, nil
}

// initializeToken performs the first API call to get the starting token.
// The messages returned in this call are NOT sent to the message channel,
// effectively discarding the chat history.
func (it *youtubeIterator) initializeToken() error {
	it.logger.Debug("Getting initial page token to start from NOW...")

	// Call without a page token to get the history and the NEXT token.
	call := it.service.LiveChatMessages.List(it.liveChatID, []string{"snippet"}).
		Context(it.ctx)

	// NOTE: We deliberately do not request authorDetails here to slightly reduce the response size
	// for a request whose data we intend to discard.

	resp, err := call.Do()
	if err != nil {
		return fmt.Errorf("initial API call failed: %w", err)
	}

	// This is the key step: Capture the token but ignore the items in 'resp.Items'.
	it.pageToken = resp.NextPageToken

	// Update polling delay based on API recommendation
	if resp.PollingIntervalMillis > 0 {
		it.pollingDelay = time.Duration(resp.PollingIntervalMillis) * time.Millisecond
	}

	it.logger.Debug(fmt.Sprintf("Initialization complete. Starting token: %s", it.pageToken))
	return nil
}

// Next blocks and returns the next message from the channel.
func (it *youtubeIterator) Next() (*ChatMessage, bool) {
	select {
	case msg, ok := <-it.messageChan:
		return msg, ok
	case <-it.ctx.Done():
		return nil, false
	}
}

// startPolling is the goroutine that continually polls the YouTube API.
func (it *youtubeIterator) startPolling() {
	defer close(it.messageChan)
	it.logger.Debug(fmt.Sprintf("Starting chat polling for ID: %s", it.liveChatID))

	for {
		select {
		case <-it.ctx.Done():
			it.logger.Info("Polling context cancelled")
			return
		case <-time.After(it.pollingDelay):
			// Proceed with the poll after the recommended delay
		}

		// 1. Call the API
		resp, err := it.service.LiveChatMessages.List(it.liveChatID, []string{"snippet", "authorDetails"}).
			PageToken(it.pageToken).
			Context(it.ctx).
			Do()

		if err != nil {
			it.logger.Error(fmt.Sprintf("Error polling API: %v. Retrying in %s.", err, it.pollingDelay))
			continue
		}

		// 2. Process and Send Messages to Channel
		for _, msg := range resp.Items {
			chatMsg := convertToChatMessage(msg)
			select {
			case it.messageChan <- chatMsg:
				// Message sent successfully
			case <-it.ctx.Done():
				// Context was cancelled while waiting to send
				return
			}
		}

		// 3. Update Polling Parameters
		it.pageToken = resp.NextPageToken
		if resp.PollingIntervalMillis > 0 {
			// Update delay based on API recommendation
			it.pollingDelay = time.Duration(resp.PollingIntervalMillis) * time.Millisecond
		}
	}
}

// getLiveChatID fetches the activeLiveChatId for a given video ID.
func (c *youtubeAPIClient) getLiveChatID(videoID string) (string, error) {
	// (Implementation remains the same as previous example)
	call := c.service.Videos.List([]string{"liveStreamingDetails"}).Id(videoID)
	response, err := call.Do()
	if err != nil {
		return "", fmt.Errorf("error calling videos.list: %w", err)
	}
	if len(response.Items) == 0 || response.Items[0].LiveStreamingDetails == nil || response.Items[0].LiveStreamingDetails.ActiveLiveChatId == "" {
		return "", fmt.Errorf("video %s not a current live stream with an active chat", videoID)
	}
	return response.Items[0].LiveStreamingDetails.ActiveLiveChatId, nil
}

// convertToChatMessage maps the complex YouTube struct to your simple, public struct.
func convertToChatMessage(ym *youtube.LiveChatMessage) *ChatMessage {
	timestamp, _ := time.Parse(time.RFC3339, ym.Snippet.PublishedAt)

	// Default to display message, which works for standard and Super Chat text
	message := ym.Snippet.DisplayMessage

	// Special handling for SuperChat to include price
	if ym.Snippet.Type == "superChatEvent" && ym.Snippet.SuperChatDetails != nil {
		details := ym.Snippet.SuperChatDetails
		// Concatenate currency and amount for a meaningful message field
		message = fmt.Sprintf("[%s %s] %s", details.Currency, formatMoney(details.AmountMicros), message)
	}

	return &ChatMessage{
		ID:        ym.Id,
		Author:    ym.AuthorDetails.DisplayName,
		Message:   message,
		Type:      ym.Snippet.Type,
		Timestamp: timestamp,
	}
}

// Helper to convert micros to standard decimal format (e.g., 1000000 -> 1.00)
func formatMoney(micros uint64) string {
	// Assuming 6 decimal places for micros
	f := float64(micros) / 1000000.0
	return strconv.FormatFloat(f, 'f', 2, 64)
}
