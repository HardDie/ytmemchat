package config

type Alerts struct {
	Enabled          bool
	Token            rune
	HTTPPort         string
	MediaPath        string
	CommandsFilePath string
}

func alertsConfig() Alerts {
	return Alerts{
		Enabled:          getEnvAsBool("ALERTS_ENABLED"),
		Token:            getEnvAsRune("ALERTS_TOKEN"),
		HTTPPort:         getEnv("ALERTS_HTTP_PORT"),
		MediaPath:        getEnv("ALERTS_MEDIA_PATH"),
		CommandsFilePath: getEnv("ALERTS_COMMANDS_FILE_PATH"),
	}
}
