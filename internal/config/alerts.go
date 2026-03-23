package config

// Alerts contains settings for trigger-based notifications and media paths.
type Alerts struct {
	Enabled          bool
	Token            rune
	MediaPath        string
	CommandsFilePath string
}

func alertsConfig() Alerts {
	return Alerts{
		Enabled:          getEnvAsBool("ALERTS_ENABLED"),
		Token:            getEnvAsRune("ALERTS_TOKEN"),
		MediaPath:        getEnv("ALERTS_MEDIA_PATH"),
		CommandsFilePath: getEnv("ALERTS_COMMANDS_FILE_PATH"),
	}
}
