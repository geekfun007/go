package handler

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/example/microservices-project/api-gateway/pkg/response"
)

// Ping 健康检查
func Ping(ctx context.Context, c *app.RequestContext) {
	response.Success(c, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "api-gateway",
	})
}

// Health 健康检查（更详细）
func Health(ctx context.Context, c *app.RequestContext) {
	response.Success(c, map[string]interface{}{
		"status": "healthy",
		"services": map[string]string{
			"api-gateway":  "ok",
			"user-service": "ok", // TODO: 实际检查用户服务状态
		},
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
