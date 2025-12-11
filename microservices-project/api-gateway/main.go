package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/example/microservices-project/api-gateway/biz/client"
	"github.com/example/microservices-project/api-gateway/biz/middleware"
	"github.com/example/microservices-project/api-gateway/biz/router"
	"github.com/example/microservices-project/api-gateway/conf"
)

func main() {
	// 加载配置
	cfg, err := conf.LoadConfig("conf/config.yaml")
	if err != nil {
		log.Printf("Failed to load config: %v, using default config", err)
		cfg = conf.GetConfig()
	}

	// 初始化用户服务客户端
	client.InitUserClient(cfg.Etcd.Endpoints)
	hlog.Infof("User client initialized with etcd endpoints: %v", cfg.Etcd.Endpoints)

	// 创建 Hertz 服务器
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)),
		server.WithReadTimeout(time.Duration(cfg.Server.ReadTimeout)*time.Second),
		server.WithWriteTimeout(time.Duration(cfg.Server.WriteTimeout)*time.Second),
	)

	// 注册全局中间件
	h.Use(middleware.Recovery()) // 异常恢复（放在最前面）
	h.Use(middleware.Logger())   // 日志记录
	h.Use(middleware.CORS())     // CORS 跨域

	// 注册路由
	router.RegisterRoutes(h)

	hlog.Infof("API Gateway starting on %s:%d", cfg.Server.Address, cfg.Server.Port)
	hlog.Info("Routes registered:")
	hlog.Info("  - GET  /ping")
	hlog.Info("  - GET  /health")
	hlog.Info("  - POST /api/v1/auth/register")
	hlog.Info("  - POST /api/v1/auth/login")
	hlog.Info("  - GET  /api/v1/users")
	hlog.Info("  - GET  /api/v1/users/:id")
	hlog.Info("  - GET  /api/v1/users/me (auth required)")
	hlog.Info("  - PUT  /api/v1/users/:id (auth required)")
	hlog.Info("  - DELETE /api/v1/users/:id (auth required)")

	h.Spin()
}
