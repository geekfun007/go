package conf

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Etcd   EtcdConfig   `yaml:"etcd"`
	JWT    JWTConfig    `yaml:"jwt"`
	Log    LogConfig    `yaml:"log"`
	CORS   CORSConfig   `yaml:"cors"`
}

type ServerConfig struct {
	Address      string `yaml:"address"`
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

type EtcdConfig struct {
	Endpoints []string `yaml:"endpoints"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type LogConfig struct {
	Level    string `yaml:"level"`
	Filename string `yaml:"filename"`
}

type CORSConfig struct {
	AllowOrigins []string `yaml:"allow_origins"`
	AllowMethods []string `yaml:"allow_methods"`
	AllowHeaders []string `yaml:"allow_headers"`
}

var GlobalConfig *Config

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	GlobalConfig = &config
	return &config, nil
}

func GetConfig() *Config {
	if GlobalConfig == nil {
		// 使用默认配置
		GlobalConfig = &Config{
			Server: ServerConfig{
				Address:      "0.0.0.0",
				Port:         8080,
				ReadTimeout:  5,
				WriteTimeout: 5,
			},
			Etcd: EtcdConfig{
				Endpoints: []string{"127.0.0.1:2379"},
			},
			JWT: JWTConfig{
				Secret:      "your-super-secret-key-change-in-production",
				ExpireHours: 24,
			},
			Log: LogConfig{
				Level:    "info",
				Filename: "logs/api-gateway.log",
			},
			CORS: CORSConfig{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders: []string{"Content-Type", "Authorization"},
			},
		}
	}
	return GlobalConfig
}
