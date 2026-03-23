// Package server provides a real-time web interface and WebSocket hub
// for the OBS overlay. It serves the HTML/CSS/JS frontend and broadcasts
// media events (alerts, TTS) to all connected browser sources.
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/HardDie/ytmemchat/pkg/logger"
)

// Server represents the HTTP and WebSocket coordinator.
type Server struct {
	cfg    Config
	mux    *http.ServeMux
	srv    http.Server
	logger *slog.Logger
}

// New initializes the server, registers core routes (/, /ws, /favicon),
// and starts the background broadcaster goroutine.
func New(cfg Config) *Server {
	mux := http.NewServeMux()

	srv := &Server{
		cfg: cfg,
		mux: mux,

		srv: http.Server{
			Addr:              cfg.Port,
			ReadHeaderTimeout: 5 * time.Second,
		},
		logger: logger.Logger.With(slog.String(logger.LogService, "server")),
	}

	// HTTP/HTML Route: Handles the root path (/)
	mux.HandleFunc("/", htmlHandler)
	// WebSocket Route: Handles real-time client connections
	mux.HandleFunc("/ws", srv.wsHandler)
	// Favicon path
	mux.HandleFunc("/favicon.ico", faviconHandler)

	// Start the broadcaster in a goroutine
	go srv.broadcaster()

	return srv
}

func (s *Server) RegisterHandle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *Server) RegisterHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *Server) GetBroadcast() chan WebsocketPayload {
	return broadcast
}

// Run starts the HTTP server. This is a blocking call.
func (s *Server) Run() error {
	s.logger.Info(fmt.Sprintf("serving http server on %s", s.cfg.Port))
	s.srv.Handler = s.mux
	err := s.srv.ListenAndServe()
	if err != nil {
		s.logger.Error(
			"ListenAndServe finished with error",
			slog.String(logger.LogValueError, err.Error()),
		)
		return fmt.Errorf("srv.ListenAndServe(): %w", err)
	}
	return nil
}

// GracefulShutdown ensures all connections are closed cleanly during app termination.
func (s *Server) GracefulShutdown(_ error) {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownRelease()

	s.logger.Info("graceful shutdown started...")
	s.srv.SetKeepAlivesEnabled(false)
	err := s.srv.Shutdown(shutdownCtx)
	if err != nil {
		s.logger.Error(
			"graceful shutdown finished with error",
			slog.String(logger.LogValueError, err.Error()),
		)
		return
	}
	s.logger.Info("graceful shutdown finished!")
}
