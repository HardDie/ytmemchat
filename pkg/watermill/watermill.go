package watermill

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/HardDie/ytmemchat/pkg/logger"
)

type Watermill struct {
	router *message.Router
	pubsub *gochannel.GoChannel
	logger *slog.Logger
}

func New(cfg Config) (*Watermill, error) {
	lg := watermill.NewStdLogger(cfg.Debug, cfg.Trace)
	pubSub := gochannel.NewGoChannel(gochannel.Config{
		OutputChannelBuffer:            0,
		Persistent:                     false,
		BlockPublishUntilSubscriberAck: false,
	}, lg)

	router, err := message.NewRouter(message.RouterConfig{
		CloseTimeout: 0,
	}, lg)
	if err != nil {
		return nil, fmt.Errorf("failed to create router: %w", err)
	}

	return &Watermill{
		router: router,
		pubsub: pubSub,
		logger: logger.Logger.With(slog.String(logger.LogService, "watermill")),
	}, nil
}

func (w *Watermill) RegisterHandler(name string, topic Topic, handler message.NoPublishHandlerFunc) {
	w.router.AddConsumerHandler(name, string(topic), w.pubsub, handler)
}

func (w *Watermill) Run() error {
	return w.router.Run(context.Background())
}

func (w *Watermill) GracefulShutdown(_ error) {
	w.logger.Info("graceful shutdown started...")
	if err := w.router.Close(); err != nil {
		w.logger.Error(
			"graceful shutdown finished with error",
			slog.String(logger.LogValueError, err.Error()),
		)
		return
	}
	w.logger.Info("graceful shutdown finished!")
}

func (w *Watermill) Publish(ctx context.Context, topic Topic, payload any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %w", err)
	}
	msg := message.NewMessage(watermill.NewUUID(), payloadBytes)
	msg.SetContext(ctx)
	err = w.pubsub.Publish(string(topic), msg)
	if err != nil {
		return fmt.Errorf("pubsub.Publish(%s): %w", topic, err)
	}
	return nil
}

func (w *Watermill) IsReady() {
	<-w.router.Running()
}
