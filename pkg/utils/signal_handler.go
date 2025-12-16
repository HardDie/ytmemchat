package utils

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/HardDie/ytmemchat/pkg/logger"
)

type SignalHandler struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger *slog.Logger
}

func NewSignalHandler(ctx context.Context) (context.Context, *SignalHandler) {
	ctx, cancel := context.WithCancel(ctx)
	return ctx, &SignalHandler{
		ctx:    ctx,
		cancel: cancel,
		logger: logger.Logger.With(slog.String(logger.LogService, "pkg/signal_handler")),
	}
}

func (h *SignalHandler) Run() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-c:
		h.logger.Debug(
			"got signal",
			slog.String("signal", sig.String()),
		)
		return fmt.Errorf("received signal: %v", sig)
	case <-h.ctx.Done():
		h.logger.Debug("context done")
		return h.ctx.Err()
	}
}

func (h *SignalHandler) GracefulShutdown(_ error) {
	h.logger.Info("graceful shutdown started...")
	h.cancel()
	h.logger.Info("graceful shutdown finished!")
}
