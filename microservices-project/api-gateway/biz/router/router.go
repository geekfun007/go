package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/example/microservices-project/api-gateway/biz/handler"
	"github.com/example/microservices-project/api-gateway/biz/middleware"
)

func RegisterRoutes(h *server.Hertz) {
	// 健康检查（无需认证）
	h.GET("/ping", handler.Ping)
	h.GET("/health", handler.Health)

	// API v1 分组
	v1 := h.Group("/api/v1")
	{
		// 认证相关（无需登录）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}

		// 用户管理（需要登录）
		users := v1.Group("/users")
		{
			// 公开接口
			users.GET("", handler.ListUsers)           // 列出用户
			users.GET("/:id", handler.GetUser)         // 获取用户信息

			// 需要认证的接口
			usersAuth := users.Group("", middleware.Auth())
			{
				usersAuth.GET("/me", handler.GetCurrentUser)      // 获取当前用户信息
				usersAuth.PUT("/:id", handler.UpdateUser)         // 更新用户信息
				usersAuth.DELETE("/:id", handler.DeleteUser)      // 删除用户
			}
		}
	}
}
