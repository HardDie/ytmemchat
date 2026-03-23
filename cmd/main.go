package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/HardDie/ytmemchat/internal/webhook"
	"github.com/HardDie/ytmemchat/pkg/utils"
	"github.com/oklog/run"

	"github.com/HardDie/ytmemchat/internal/alerts"
	clientYoutube "github.com/HardDie/ytmemchat/internal/clients/youtube"
	clientYoutubeV1 "github.com/HardDie/ytmemchat/internal/clients/youtubev1"
	"github.com/HardDie/ytmemchat/internal/config"
	"github.com/HardDie/ytmemchat/internal/server"
	"github.com/HardDie/ytmemchat/internal/tts"
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
	cfg := config.Get()
	ctx, signalHandler := utils.NewSignalHandler(context.Background())

	srv := server.New(server.Config{
		Port: cfg.Server.Port,
	})

	var err error
	var yt clientYoutube.Client
	if true {
		yt, err = clientYoutube.New(cfg.Youtube.APIKey)
		if err != nil {
			logger.Error(
				"failed to create youtube client",
				slog.String(logger.LogValueError, err.Error()),
			)
			return exitFailure
		}
	} else {
		yt, err = clientYoutubeV1.New()
		if err != nil {
			logger.Error(
				"failed to create youtube client v1",
				slog.String(logger.LogValueError, err.Error()),
			)
			return exitFailure
		}
	}
	ytIt, err := yt.GetMessageIterator(ctx, cfg.Youtube.StreamID)
	if err != nil {
		logger.Error(
			"failed to run youtube iterator",
			slog.String(logger.LogValueError, err.Error()),
		)
		return exitFailure
	}

	al, err := alerts.New(alerts.Config{
		Token:            string(cfg.Alerts.Token),
		MediaPath:        cfg.Alerts.MediaPath,
		CommandsFilePath: cfg.Alerts.CommandsFilePath,
		Broadcast:        srv.GetBroadcast(),
	})
	if err != nil {
		logger.Error(
			"failed to init alert service",
			slog.String(logger.LogValueError, err.Error()),
		)
		return exitFailure
	}
	ttsService := tts.New(tts.Config{
		VoiceName: cfg.TTS.Name,
		Broadcast: srv.GetBroadcast(),
	})
	whService := webhook.New(webhook.Config{
		Broadcast: srv.GetBroadcast(),
	})

	if cfg.Alerts.Enabled {
		srv.RegisterHandle("/media/", al.GetMediaHandler())
	}
	if cfg.Webhook.Enabled {
		srv.RegisterHandleFunc("/webhook/", whService.Handle)
		srv.RegisterHandleFunc("/interrupt/", whService.InterruptHandle)
	}

	// Run all background services with graceful shutdown
	var g run.Group
	g.Add(srv.Run, srv.GracefulShutdown)
	g.Add(
		func() error {
			for {
				if ctx.Err() != nil {
					return ctx.Err()
				}

				var message *clientYoutube.ChatMessage
				var ok bool

				select {
				case message, ok = <-ytIt.GetChan():
				case message, _ = <-whService.GetChan():
					// always true
					ok = true
				case <-ctx.Done():
					ok = false
				}
				//message, ok := ytIt.Next()
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

				if cfg.Alerts.Enabled {
					if al.Alert(message.Message) {
						// Do not pronounce messages with alert token.
						continue
					}
				}

				if cfg.TTS.Enabled {
					err = ttsService.SynthesizeAudio(message.Message)
					if err != nil {
						logger.Error(
							"failed to speak message",
							slog.String(logger.LogValueError, err.Error()),
							slog.String(logger.LogMessage, message.Message),
							slog.String(logger.LogTTSName, cfg.TTS.Name),
						)
					}
				}
			}
			return nil
		},
		func(err error) {},
	)
	g.Add(signalHandler.Run, signalHandler.GracefulShutdown)

	// Working!
	if err = g.Run(); err != nil {
		logger.Error(
			"error running group",
			slog.String("error", err.Error()),
		)
		return exitFailure
	}

	return exitSuccess
}
