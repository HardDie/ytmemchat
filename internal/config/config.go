// Package config handles the application configuration by loading environment
// variables from a .env file or the system environment.
//
// It uses a fail-fast approach: if a required variable is missing or
// malformed, the application will log a fatal error and panic during initialization.
package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/HardDie/ytmemchat/pkg/logger"
)

// Config holds the centralized configuration for the entire application.
type Config struct {
	Alerts  Alerts
	Server  Server
	TTS     TTS
	Webhook Webhook
	Youtube Youtube
}

// Get initializes and returns the Config struct.
// It attempts to load a .env file on startup. If the file is missing,
// it continues to check system environment variables.
func Get() Config {
	if err := godotenv.Load(); err != nil {
		if check := os.IsNotExist(err); !check {
			logger.Error(
				"failed to load env vars",
				slog.String(logger.LogValueError, err.Error()),
			)
			//nolint:govet // it's okay for configuration
			panic(nil)
		}
	}

	cfg := Config{
		Alerts:  alertsConfig(),
		Server:  serverConfig(),
		TTS:     ttsConfig(),
		Webhook: webhookConfig(),
		Youtube: youtubeConfig(),
	}
	return cfg
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		logger.Error(
			"env value not found",
			slog.String("key", key),
		)
		//nolint:govet // it's okay for configuration
		panic(nil)
	}
	return value
}

func getEnvAsDuration(key string) time.Duration {
	value := getEnv(key)
	v, err := time.ParseDuration(value)
	if err != nil {
		logger.Error(
			"env value invalid duration",
			slog.String("key", key),
		)
		//nolint:govet // it's okay for configuration
		panic(nil)
	}
	return v
}

func getEnvAsInt(key string) int {
	value := getEnv(key)
	v, e := strconv.Atoi(value)
	if e != nil {
		logger.Error(
			"env value invalid int",
			slog.String("key", key),
		)
		//nolint:govet // it's okay for configuration
		panic(nil)
	}
	return v
}

func getEnvAsBool(key string) bool {
	value := getEnv(key)
	switch {
	case strings.EqualFold(value, "true"):
		return true
	case strings.EqualFold(value, "false"):
		return false
	default:
		logger.Error(
			"env value invalid bool",
			slog.String("key", key),
		)
		//nolint:govet // it's okay for configuration
		panic(nil)
	}
}

func getEnvAsRune(key string) rune {
	value := []rune(getEnv(key))
	if len(value) != 1 {
		logger.Error(
			"env value invalid rune",
			slog.String("key", key),
		)
		//nolint:govet // it's okay for configuration
		panic(nil)
	}
	return value[0]
}
