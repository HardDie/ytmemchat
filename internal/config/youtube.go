package config

import "os"

type Youtube struct {
	APIKey   string
	StreamID string
}

func youtubeConfig() Youtube {
	return Youtube{
		APIKey:   os.Getenv("YOUTUBE_API_KEY"),
		StreamID: os.Getenv("YOUTUBE_STREAM_ID"),
	}
}
