package config

import "os"

type Config struct {
	HTTP     HTTPConfig
	RPC      RPCConfig
	Etcd     EtcdConfig
	Registry RegistryConfig
}

type HTTPConfig struct {
	Host string
	Port string
}

type RPCConfig struct {
	Name string // RPC 服务名称（用于服务发现）
}

type EtcdConfig struct {
	Endpoints []string
}

type RegistryConfig struct {
	Enable bool
}

func GetConfig() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Host: getEnv("HTTP_HOST", "0.0.0.0"),
			Port: getEnv("HTTP_PORT", "8080"),
		},
		RPC: RPCConfig{
			Name: getEnv("RPC_SERVICE_NAME", "user_service"),
		},
		Etcd: EtcdConfig{
			Endpoints: []string{getEnv("ETCD_ENDPOINTS", "localhost:2379")},
		},
		Registry: RegistryConfig{
			Enable: getEnv("REGISTRY_ENABLE", "true") == "true",
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
