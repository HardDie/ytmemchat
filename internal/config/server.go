package config

// Server defines the network configuration for the application.
type Server struct {
	Port string
}

func serverConfig() Server {
	return Server{
		Port: getEnv("SERVER_PORT"),
	}
}
