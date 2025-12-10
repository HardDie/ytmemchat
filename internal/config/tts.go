package config

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
