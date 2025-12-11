package conf

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Etcd     EtcdConfig     `yaml:"etcd"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
	Address     string `yaml:"address"`
	Port        int    `yaml:"port"`
	ServiceName string `yaml:"service_name"`
}

type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
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
				Address:     "0.0.0.0",
				Port:        8888,
				ServiceName: "user-service",
			},
			Database: DatabaseConfig{
				Host:            "127.0.0.1",
				Port:            3306,
				User:            "root",
				Password:        "root123",
				Database:        "userdb",
				MaxIdleConns:    10,
				MaxOpenConns:    100,
				ConnMaxLifetime: 3600,
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
				Filename: "logs/user-service.log",
			},
		}
	}
	return GlobalConfig
}
