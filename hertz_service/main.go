package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/example/hertz-kitex-demo/hertz_service/client"
	"github.com/example/hertz-kitex-demo/hertz_service/config"
	"github.com/example/hertz-kitex-demo/hertz_service/handler"
)

func main() {
	// 加载配置
	cfg := config.GetConfig()

	// 初始化 RPC 客户端
	if err := client.InitUserClient(&cfg.RPC); err != nil {
		log.Fatalf("Failed to initialize RPC client: %v", err)
	}

	// 创建 Hertz HTTP 服务器
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)),
	)

	// 创建 handler
	userHandler := handler.NewUserHandler()

	// 注册路由
	registerRoutes(h, userHandler)

	log.Printf("Hertz HTTP Server starting at %s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	// 启动服务
	h.Spin()
}

func registerRoutes(h *server.Hertz, userHandler *handler.UserHandler) {
	// 健康检查
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{
			"message": "pong",
		})
	})

	// API 路由组
	api := h.Group("/api/v1")
	{
		// 用户相关路由
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)       // 创建用户
			users.GET("/:id", userHandler.GetUser)       // 获取用户
			users.PUT("/:id", userHandler.UpdateUser)    // 更新用户
			users.DELETE("/:id", userHandler.DeleteUser) // 删除用户
			users.GET("", userHandler.ListUsers)         // 获取用户列表
		}
	}
}
