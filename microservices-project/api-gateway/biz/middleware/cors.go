package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/example/microservices-project/api-gateway/conf"
)

func CORS() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		cfg := conf.GetConfig()

		// 设置 CORS 头
		c.Header("Access-Control-Allow-Origin", strings.Join(cfg.CORS.AllowOrigins, ","))
		c.Header("Access-Control-Allow-Methods", strings.Join(cfg.CORS.AllowMethods, ","))
		c.Header("Access-Control-Allow-Headers", strings.Join(cfg.CORS.AllowHeaders, ","))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// 处理 OPTIONS 预检请求
		if string(c.Method()) == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next(ctx)
	}
}
