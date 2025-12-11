package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Logger() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := string(c.Path())
		method := string(c.Method())

		// 处理请求
		c.Next(ctx)

		// 记录日志
		duration := time.Since(start)
		statusCode := c.Response.StatusCode()

		hlog.Infof("[%s] %s %s | %d | %v | %s",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			statusCode,
			duration,
			c.ClientIP(),
		)
	}
}
