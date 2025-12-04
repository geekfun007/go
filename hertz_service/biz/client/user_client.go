package client

import (
	"fmt"
	"log"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/example/hertz-kitex-demo/hertz_service/biz/config"
	"github.com/example/hertz-kitex-demo/kitex_gen/user/userservice"
)

var userClient userservice.Client

// InitUserClient 初始化用户服务客户端
func InitUserClient(cfg *config.Config) error {
	var opts []client.Option

	// 如果启用服务发现
	if cfg.Registry.Enable {
		log.Printf("Service discovery enabled with etcd: %v", cfg.Etcd.Endpoints)
		
		// 创建 etcd resolver
		r, err := etcd.NewEtcdResolver(cfg.Etcd.Endpoints)
		if err != nil {
			return fmt.Errorf("failed to create etcd resolver: %w", err)
		}

		opts = append(opts, client.WithResolver(r))
		log.Printf("Connecting to RPC service via service discovery: %s", cfg.RPC.Name)
	} else {
		// 如果不使用服务发现，直连模式（需要指定地址）
		log.Println("Service discovery disabled, using direct connection")
	}

	var err error
	userClient, err = userservice.NewClient(cfg.RPC.Name, opts...)
	if err != nil {
		return fmt.Errorf("failed to create user client: %w", err)
	}

	log.Printf("User RPC client initialized for service: %s", cfg.RPC.Name)
	return nil
}

// GetUserClient 获取用户服务客户端
func GetUserClient() userservice.Client {
	return userClient
}
