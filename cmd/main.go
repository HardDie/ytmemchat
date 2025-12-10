package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	clientYoutube "github.com/HardDie/ytmemchat/internal/clients/youtube"
	"github.com/HardDie/ytmemchat/internal/config"
	"github.com/HardDie/ytmemchat/pkg/logger"
)

const (
	exitSuccess = 0
	exitFailure = 1
)

func main() {
	os.Exit(gracefulMain())
}

func gracefulMain() int {
	ctx := context.Background()
	cfg := config.Get()

	yt, err := clientYoutube.New(cfg.Youtube.APIKey)
	if err != nil {
		logger.Error(
			"failed to create youtube client",
			slog.String(logger.LogValueError, err.Error()),
		)
		return exitFailure
	}
	ytIt, err := yt.GetMessageIterator(ctx, cfg.Youtube.StreamID)
	if err != nil {
		logger.Error(
			"failed to run youtube iterator",
			slog.String(logger.LogValueError, err.Error()),
		)
		return exitFailure
	}

	for {
		message, ok := ytIt.Next()
		if !ok {
			logger.Error(
				"Youtube iterator closed. Exit application.",
			)
			break
		}
		logger.Debug(fmt.Sprintf(
			"[%s | %s] %s: %s",
			message.Timestamp.Format("15:04:05"),
			message.Type,
			message.Author,
			message.Message,
		))
	}

	return exitSuccess
}
