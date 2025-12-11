package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/example/microservices-project/api-gateway/conf"
	jwtutil "github.com/example/microservices-project/api-gateway/pkg/jwt"
	"github.com/example/microservices-project/api-gateway/pkg/response"
)

func Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 获取 Authorization header
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证 token
		cfg := conf.GetConfig()
		userID, err := jwtutil.ValidateToken(tokenString, cfg.JWT.Secret)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		// 将用户 ID 存储到上下文
		c.Set("user_id", userID)

		c.Next(ctx)
	}
}

// GetUserID 从上下文中获取用户 ID
func GetUserID(c *app.RequestContext) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}
