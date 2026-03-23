package config

type Webhook struct {
	Enabled bool
}

func webhookConfig() Webhook {
	return Webhook{
		Enabled: getEnvAsBool("WEBHOOK_ENABLED"),
	}
}
