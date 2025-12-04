package client

import (
	"fmt"
	"log"

	"github.com/cloudwego/kitex/client"
	"github.com/example/hertz-kitex-demo/hertz_service/config"
	"github.com/example/hertz-kitex-demo/kitex_service/kitex_gen/user/userservice"
)

var userClient userservice.Client

// InitUserClient 初始化用户服务客户端
func InitUserClient(cfg *config.RPCConfig) error {
	var err error
	rpcAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	
	userClient, err = userservice.NewClient(
		"user_service",
		client.WithHostPorts(rpcAddr),
	)
	if err != nil {
		return fmt.Errorf("failed to create user client: %w", err)
	}

	log.Printf("User RPC client initialized, connecting to %s", rpcAddr)
	return nil
}

// GetUserClient 获取用户服务客户端
func GetUserClient() userservice.Client {
	return userClient
}
