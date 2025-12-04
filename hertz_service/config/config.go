package config

import "os"

type Config struct {
	HTTP HTTPConfig
	RPC  RPCConfig
}

type HTTPConfig struct {
	Host string
	Port string
}

type RPCConfig struct {
	Host string
	Port string
}

func GetConfig() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Host: getEnv("HTTP_HOST", "0.0.0.0"),
			Port: getEnv("HTTP_PORT", "8080"),
		},
		RPC: RPCConfig{
			Host: getEnv("RPC_HOST", "localhost"),
			Port: getEnv("RPC_PORT", "8888"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
