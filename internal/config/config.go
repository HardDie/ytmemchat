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

type Config struct {
	TTS     TTS
	Youtube Youtube
}

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
		TTS:     ttsConfig(),
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
