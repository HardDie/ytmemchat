package config

type Server struct {
	Port string
}

func serverConfig() Server {
	return Server{
		Port: getEnv("SERVER_PORT"),
	}
}
