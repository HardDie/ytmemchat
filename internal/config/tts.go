package config

// TTS defines settings for the Text-to-Speech engine.
type TTS struct {
	Enabled bool
	Name    string
}

func ttsConfig() TTS {
	return TTS{
		Enabled: getEnvAsBool("TTS_ENABLED"),
		Name:    getEnv("TTS_NAME"),
	}
}
