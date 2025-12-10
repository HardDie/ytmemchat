package logger

import (
	"log/slog"
	"os"
)

const (
	LogService    = "service"
	LogValueError = "error"
	LogMessage    = "message"
	LogTTSName    = "tts_name"
)

var (
	Logger = slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource:   true,
				Level:       slog.LevelDebug,
				ReplaceAttr: nil,
			},
		),
	)
	Debug = Logger.Debug
	Info  = Logger.Info
	Warn  = Logger.Warn
	Error = Logger.Error
)
