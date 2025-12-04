package main

import (
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/example/hertz-kitex-demo/kitex_gen/user/userservice"
	"github.com/example/hertz-kitex-demo/kitex_service/biz/handler"
	"github.com/example/hertz-kitex-demo/kitex_service/config"
	"github.com/example/hertz-kitex-demo/kitex_service/dal"
)

func main() {
	// 加载配置
	cfg := config.GetConfig()

	// 初始化数据库
	if err := dal.InitDB(&cfg.MySQL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// 创建 RPC 服务地址
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", cfg.RPC.Host, cfg.RPC.Port))
	if err != nil {
		log.Fatalf("Failed to resolve address: %v", err)
	}

	// 服务选项
	var opts []server.Option
	opts = append(opts, server.WithServiceAddr(addr))

	// 如果启用服务注册
	if cfg.Registry.Enable {
		log.Printf("Service registry enabled with etcd: %v", cfg.Etcd.Endpoints)
		
		// 创建 etcd registry
		r, err := etcd.NewEtcdRegistry(cfg.Etcd.Endpoints)
		if err != nil {
			log.Fatalf("Failed to create etcd registry: %v", err)
		}

		// 添加 registry 选项
		opts = append(opts, server.WithRegistry(r))
		opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: cfg.RPC.Name,
		}))
	}

	// 创建服务实例
	svr := userservice.NewServer(
		handler.NewUserServiceImpl(),
		opts...,
	)

	log.Printf("Kitex RPC Server starting at %s:%s", cfg.RPC.Host, cfg.RPC.Port)
	if cfg.Registry.Enable {
		log.Printf("Service registered as: %s", cfg.RPC.Name)
	}

	// 启动服务
	err = svr.Run()
	if err != nil {
		log.Fatalf("Server run failed: %v", err)
	}
}
