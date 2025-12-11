package middleware

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/example/microservices-project/api-gateway/pkg/response"
)

func Recovery() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				hlog.Errorf("Panic recovered: %v", err)

				// 返回 500 错误
				response.InternalServerError(c, fmt.Sprintf("internal server error: %v", err))
				c.Abort()
			}
		}()

		c.Next(ctx)
	}
}
