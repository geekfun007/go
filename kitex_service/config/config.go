package config

import (
	"fmt"
	"os"
)

type Config struct {
	MySQL    MySQLConfig
	RPC      RPCConfig
	Etcd     EtcdConfig
	Registry RegistryConfig
}

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type RPCConfig struct {
	Host string
	Port string
	Name string // 服务名称
}

type EtcdConfig struct {
	Endpoints []string
}

type RegistryConfig struct {
	Enable bool
}

func GetConfig() *Config {
	return &Config{
		MySQL: MySQLConfig{
			Host:     getEnv("MYSQL_HOST", "localhost"),
			Port:     getEnv("MYSQL_PORT", "3306"),
			User:     getEnv("MYSQL_USER", "root"),
			Password: getEnv("MYSQL_PASSWORD", "root"),
			Database: getEnv("MYSQL_DATABASE", "user_db"),
		},
		RPC: RPCConfig{
			Host: getEnv("RPC_HOST", "0.0.0.0"),
			Port: getEnv("RPC_PORT", "8888"),
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

func (c *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
