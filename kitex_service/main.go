package main

import (
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/example/hertz-kitex-demo/kitex_service/config"
	"github.com/example/hertz-kitex-demo/kitex_service/dal"
	user "github.com/example/hertz-kitex-demo/kitex_service/kitex_gen/user/userservice"
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

	// 创建服务实例
	svr := user.NewServer(
		NewUserServiceImpl(),
		server.WithServiceAddr(addr),
	)

	log.Printf("Kitex RPC Server starting at %s:%s", cfg.RPC.Host, cfg.RPC.Port)

	// 启动服务
	err = svr.Run()
	if err != nil {
		log.Fatalf("Server run failed: %v", err)
	}
}
