package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/HardDie/ytmemchat/pkg/logger"
)

type Server struct {
	cfg    Config
	mux    *http.ServeMux
	srv    http.Server
	logger *slog.Logger
}

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
