package config

type Youtube struct {
	APIKey   string
	StreamID string
}

func youtubeConfig() Youtube {
	return Youtube{
		APIKey:   getEnv("YOUTUBE_API_KEY"),
		StreamID: getEnv("YOUTUBE_STREAM_ID"),
	}
}
