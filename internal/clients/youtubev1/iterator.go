package youtubev1

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	clientYoutube "github.com/HardDie/ytmemchat/internal/clients/youtube"
	"github.com/HardDie/ytmemchat/pkg/logger"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

type youtubeIterator struct {
	ctx          context.Context
	liveVideoID  string
	pageToken    string
	pollingDelay time.Duration
	messageChan  chan *clientYoutube.ChatMessage
	logger       *slog.Logger
}

func New() (clientYoutube.Client, error) {
	return &youtubeAPIClient{
		logger: logger.Logger.With(slog.String(logger.LogService, "YouTube Client v1")),
	}, nil
}

type youtubeAPIClient struct {
	logger *slog.Logger
}

func (c *youtubeAPIClient) GetMessageIterator(ctx context.Context, liveVideoID string) (clientYoutube.MessageIterator, error) {
	it := &youtubeIterator{
		ctx:          ctx,
		liveVideoID:  liveVideoID,
		pageToken:    "",              // Start with an empty token
		pollingDelay: 5 * time.Second, // Default initial delay
		messageChan:  make(chan *clientYoutube.ChatMessage, 10),
		logger:       c.logger,
	}

	// Perform the initial setup poll to discard old messages and get the starting token.
	if err := it.initializeToken(liveVideoID); err != nil {
		return nil, fmt.Errorf("failed to initialize chat token: %w", err)
	}

	// Start the background goroutine to handle polling and channel population
	go it.startPolling()

	return it, nil
}

func (it *youtubeIterator) initializeToken(liveVideoID string) error {
	it.logger.Debug("Getting initial page token to start from NOW...")

	client := resty.New()

	// 1. Fetch initial HTML
	resp, err := client.R().
		SetHeader("User-Agent", "Firefox/99").
		Get(fmt.Sprintf("https://www.youtube.com/live_chat?v=%s", liveVideoID))
	if err != nil {
		log.Fatal(err)
	}

	// 2. Use goquery (HTML Parser) to find the script tags
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		log.Fatal(err)
	}

	// Loop through scripts to find the one containing ytInitialData
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		scriptContent := s.Text()

		// Look for the continuation token
		if strings.Contains(scriptContent, "liveChatRenderer") {
			parts := strings.Split(scriptContent, `"continuation":"`)
			if len(parts) > 1 {
				it.pageToken = strings.Split(parts[1], `"`)[0]
			}
		}
	})

	if it.pageToken == "" {
		return fmt.Errorf("failed to extract session data using parser")
	}

	it.logger.Debug(fmt.Sprintf("Initialization complete. Starting token: %s", it.pageToken))
	return nil
}

func (it *youtubeIterator) Next() (*clientYoutube.ChatMessage, bool) {
	select {
	case msg, ok := <-it.messageChan:
		return msg, ok
	case <-it.ctx.Done():
		return nil, false
	}
}

func (it *youtubeIterator) startPolling() {
	defer close(it.messageChan)
	it.logger.Debug(fmt.Sprintf("Starting chat polling for ID: %s", it.liveVideoID))

	client := resty.New()

	for {
		select {
		case <-it.ctx.Done():
			it.logger.Info("Polling context cancelled")
			return
		case <-time.After(it.pollingDelay):
			// Proceed with the poll after the recommended delay
		}

		var result YTInternalResponse
		payload := map[string]interface{}{
			"context": map[string]interface{}{
				"client": map[string]interface{}{"clientName": "WEB", "clientVersion": "2.9999099"},
			},
			"continuation": it.pageToken,
		}

		// 1. Call the API
		resp, _ := client.R().
			SetQueryParam("prettyPrint", "false").
			SetBody(payload).
			SetResult(&result).
			Post("https://www.youtube.com/youtubei/v1/live_chat/get_live_chat")

		if resp.IsError() {
			it.logger.Error(fmt.Sprintf("Error polling API: %v. Retrying in %s.", resp.Error(), it.pollingDelay))
			continue
		}

		// 2. Process and Send Messages to Channel
		actions := result.ContinuationContents.LiveChatContinuation.Actions
		for _, action := range actions {
			chatMsg := convertToChatMessage(action.AddChatItemAction.Item.LiveChatTextMessageRenderer)
			if chatMsg.Message == "" {
				continue
			}
			select {
			case it.messageChan <- chatMsg:
				// Message sent successfully
			case <-it.ctx.Done():
				// Context was cancelled while waiting to send
				return
			}
		}

		// 3. Update Polling Parameters
		conns := result.ContinuationContents.LiveChatContinuation.Continuations
		if len(conns) > 0 {
			if conns[0].TimedContinuationData != nil {
				if conns[0].TimedContinuationData.Continuation != "" {
					it.pageToken = conns[0].TimedContinuationData.Continuation
				}
				if conns[0].TimedContinuationData.TimeoutMs > 0 {
					it.pollingDelay = time.Duration(conns[0].TimedContinuationData.TimeoutMs) * time.Millisecond
				}
			} else if conns[0].InvalidationContinuationData != nil {
				if conns[0].InvalidationContinuationData.Continuation != "" {
					it.pageToken = conns[0].InvalidationContinuationData.Continuation
				}
				if conns[0].InvalidationContinuationData.TimeoutMs > 0 {
					it.pollingDelay = time.Duration(conns[0].InvalidationContinuationData.TimeoutMs) * time.Millisecond
				}
			} else {
				it.logger.Error("continuation data has not been found")
				it.pollingDelay = 5 * time.Second
			}
		}
	}
}

func convertToChatMessage(renderer LiveChatTextMessageRenderer) *clientYoutube.ChatMessage {
	author := renderer.AuthorName.SimpleText
	msg := ""
	for _, run := range renderer.Message.Runs {
		msg += run.Text
	}

	return &clientYoutube.ChatMessage{
		Author:  author,
		Message: msg,
	}
}

type YTInternalResponse struct {
	ContinuationContents struct {
		LiveChatContinuation struct {
			Continuations []struct {
				TimedContinuationData *struct {
					Continuation string `json:"continuation"`
					TimeoutMs    int64  `json:"timeoutMs"`
				} `json:"timedContinuationData"`
				InvalidationContinuationData *struct {
					Continuation string `json:"continuation"`
					TimeoutMs    int64  `json:"timeoutMs"`
				} `json:"invalidationContinuationData"`
			} `json:"continuations"`
			Actions []struct {
				AddChatItemAction struct {
					Item struct {
						LiveChatTextMessageRenderer LiveChatTextMessageRenderer `json:"liveChatTextMessageRenderer"`
					} `json:"item"`
				} `json:"addChatItemAction"`
			} `json:"actions"`
		} `json:"liveChatContinuation"`
	} `json:"continuationContents"`
}

type LiveChatTextMessageRenderer struct {
	AuthorName struct {
		SimpleText string `json:"simpleText"`
	} `json:"authorName"`
	Message struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"message"`
}
