# Hertz - 字节跳动高性能 HTTP 框架 / ByteDance High-Performance HTTP Framework

## 1. 简介与安装 / Introduction & Installation

Hertz 是字节跳动开源的高性能 HTTP 框架，具有高易用性、高性能、高扩展性等特点。

```bash
# 安装 hz 命令行工具 / Install hz CLI
go install github.com/cloudwego/hertz/cmd/hz@latest

# 添加依赖 / Add dependencies
go get github.com/cloudwego/hertz
```

## 2. 快速开始 / Quick Start

```go
package main

import (
    "context"
    
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/common/utils"
    "github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
    h := server.Default()
    
    // 基本路由 / Basic routing
    h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
        c.JSON(consts.StatusOK, utils.H{
            "message": "pong",
        })
    })
    
    h.POST("/user", func(ctx context.Context, c *app.RequestContext) {
        c.JSON(consts.StatusOK, utils.H{
            "message": "user created",
        })
    })
    
    // 启动服务器 / Start server
    h.Spin()
}
```

## 3. 路由 / Routing

```go
package main

import (
    "context"
    
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/common/utils"
    "github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
    h := server.Default()
    
    // 路径参数 / Path parameters
    h.GET("/user/:id", func(ctx context.Context, c *app.RequestContext) {
        id := c.Param("id")
        c.JSON(consts.StatusOK, utils.H{"id": id})
    })
    
    // 通配符路由 / Wildcard routing
    h.GET("/files/*filepath", func(ctx context.Context, c *app.RequestContext) {
        filepath := c.Param("filepath")
        c.JSON(consts.StatusOK, utils.H{"path": filepath})
    })
    
    // 查询参数 / Query parameters
    h.GET("/search", func(ctx context.Context, c *app.RequestContext) {
        query := c.Query("q")
        page := c.DefaultQuery("page", "1")
        c.JSON(consts.StatusOK, utils.H{
            "query": query,
            "page":  page,
        })
    })
    
    // 路由组 / Route groups
    v1 := h.Group("/api/v1")
    {
        v1.GET("/users", listUsers)
        v1.POST("/users", createUser)
        v1.GET("/users/:id", getUser)
        v1.PUT("/users/:id", updateUser)
        v1.DELETE("/users/:id", deleteUser)
    }
    
    // 静态文件 / Static files
    h.Static("/static", "./static")
    h.StaticFile("/favicon.ico", "./static/favicon.ico")
    
    h.Spin()
}

func listUsers(ctx context.Context, c *app.RequestContext) {
    c.JSON(consts.StatusOK, utils.H{"users": []string{}})
}

func createUser(ctx context.Context, c *app.RequestContext) {
    c.JSON(consts.StatusCreated, utils.H{"message": "created"})
}

func getUser(ctx context.Context, c *app.RequestContext) {
    id := c.Param("id")
    c.JSON(consts.StatusOK, utils.H{"id": id})
}

func updateUser(ctx context.Context, c *app.RequestContext) {
    c.JSON(consts.StatusOK, utils.H{"message": "updated"})
}

func deleteUser(ctx context.Context, c *app.RequestContext) {
    c.JSON(consts.StatusOK, utils.H{"message": "deleted"})
}
```

## 4. 请求处理 / Request Handling

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/protocol/consts"
)

// 请求绑定结构体 / Request binding struct
type CreateUserReq struct {
    Name  string `json:"name" form:"name" query:"name" vd:"len($) > 0"`
    Email string `json:"email" form:"email" query:"email" vd:"email($)"`
    Age   int    `json:"age" form:"age" query:"age" vd:"$ >= 0 && $ <= 150"`
}

type User struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func main() {
    h := server.Default()
    
    // JSON 绑定 / JSON binding
    h.POST("/json", func(ctx context.Context, c *app.RequestContext) {
        var req CreateUserReq
        if err := c.BindJSON(&req); err != nil {
            c.JSON(consts.StatusBadRequest, map[string]string{
                "error": err.Error(),
            })
            return
        }
        c.JSON(consts.StatusOK, req)
    })
    
    // 表单绑定 / Form binding
    h.POST("/form", func(ctx context.Context, c *app.RequestContext) {
        var req CreateUserReq
        if err := c.Bind(&req); err != nil {
            c.JSON(consts.StatusBadRequest, map[string]string{
                "error": err.Error(),
            })
            return
        }
        c.JSON(consts.StatusOK, req)
    })
    
    // 获取请求头 / Get request headers
    h.GET("/headers", func(ctx context.Context, c *app.RequestContext) {
        contentType := c.GetHeader("Content-Type")
        userAgent := c.UserAgent()
        c.JSON(consts.StatusOK, map[string]string{
            "content-type": string(contentType),
            "user-agent":   string(userAgent),
        })
    })
    
    // 获取 Cookie / Get cookie
    h.GET("/cookie", func(ctx context.Context, c *app.RequestContext) {
        val := c.Cookie("session")
        c.JSON(consts.StatusOK, map[string]string{
            "session": string(val),
        })
    })
    
    // 文件上传 / File upload
    h.POST("/upload", func(ctx context.Context, c *app.RequestContext) {
        file, err := c.FormFile("file")
        if err != nil {
            c.JSON(consts.StatusBadRequest, map[string]string{
                "error": err.Error(),
            })
            return
        }
        
        // 保存文件 / Save file
        dst := fmt.Sprintf("./uploads/%s", file.Filename)
        c.SaveUploadedFile(file, dst)
        
        c.JSON(consts.StatusOK, map[string]string{
            "filename": file.Filename,
            "size":     fmt.Sprintf("%d", file.Size),
        })
    })
    
    // 多文件上传 / Multiple file upload
    h.POST("/uploads", func(ctx context.Context, c *app.RequestContext) {
        form, _ := c.MultipartForm()
        files := form.File["files"]
        
        for _, file := range files {
            dst := fmt.Sprintf("./uploads/%s", file.Filename)
            c.SaveUploadedFile(file, dst)
        }
        
        c.JSON(consts.StatusOK, map[string]int{
            "uploaded": len(files),
        })
    })
    
    h.Spin()
}
```

## 5. 响应处理 / Response Handling

```go
package main

import (
    "context"
    
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
    h := server.Default()
    
    // JSON 响应 / JSON response
    h.GET("/json", func(ctx context.Context, c *app.RequestContext) {
        c.JSON(consts.StatusOK, map[string]interface{}{
            "code":    0,
            "message": "success",
            "data": map[string]string{
                "name": "Alice",
            },
        })
    })
    
    // 设置响应头 / Set response headers
    h.GET("/headers", func(ctx context.Context, c *app.RequestContext) {
        c.Header("X-Custom-Header", "value")
        c.Header("Content-Type", "application/json")
        c.JSON(consts.StatusOK, map[string]string{"message": "ok"})
    })
    
    // 设置 Cookie / Set cookie
    h.GET("/set-cookie", func(ctx context.Context, c *app.RequestContext) {
        c.SetCookie("session", "abc123", 3600, "/", "localhost", false, true)
        c.JSON(consts.StatusOK, map[string]string{"message": "cookie set"})
    })
    
    // 重定向 / Redirect
    h.GET("/redirect", func(ctx context.Context, c *app.RequestContext) {
        c.Redirect(consts.StatusFound, []byte("/target"))
    })
    
    // 字符串响应 / String response
    h.GET("/string", func(ctx context.Context, c *app.RequestContext) {
        c.String(consts.StatusOK, "Hello, %s!", "World")
    })
    
    // HTML 响应 / HTML response
    h.GET("/html", func(ctx context.Context, c *app.RequestContext) {
        c.HTML(consts.StatusOK, "<h1>Hello, World!</h1>")
    })
    
    // 文件下载 / File download
    h.GET("/download", func(ctx context.Context, c *app.RequestContext) {
        c.File("./files/document.pdf")
    })
    
    // 流式响应 / Streaming response
    h.GET("/stream", func(ctx context.Context, c *app.RequestContext) {
        c.SetStatusCode(consts.StatusOK)
        c.Response.Header.Set("Content-Type", "text/event-stream")
        c.Response.Header.Set("Cache-Control", "no-cache")
        
        for i := 0; i < 5; i++ {
            c.WriteString(fmt.Sprintf("data: message %d\n\n", i))
            c.Flush()
            time.Sleep(time.Second)
        }
    })
    
    h.Spin()
}
```

## 6. 中间件 / Middleware

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/protocol/consts"
)

// 日志中间件 / Logging middleware
func LoggingMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        start := time.Now()
        path := string(c.Path())
        method := string(c.Method())
        
        c.Next(ctx)  // 处理请求
        
        latency := time.Since(start)
        status := c.Response.StatusCode()
        
        log.Printf("[%s] %s %d %v", method, path, status, latency)
    }
}

// 恢复中间件 / Recovery middleware
func RecoveryMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Panic recovered: %v", r)
                c.JSON(consts.StatusInternalServerError, map[string]string{
                    "error": "Internal Server Error",
                })
                c.Abort()
            }
        }()
        c.Next(ctx)
    }
}

// 认证中间件 / Auth middleware
func AuthMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        token := c.GetHeader("Authorization")
        if len(token) == 0 {
            c.JSON(consts.StatusUnauthorized, map[string]string{
                "error": "Unauthorized",
            })
            c.Abort()
            return
        }
        
        // 验证 token...
        // Validate token...
        
        c.Set("user_id", 123)  // 设置上下文值
        c.Next(ctx)
    }
}

// CORS 中间件 / CORS middleware
func CORSMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if string(c.Method()) == "OPTIONS" {
            c.AbortWithStatus(consts.StatusNoContent)
            return
        }
        
        c.Next(ctx)
    }
}

func main() {
    h := server.Default()
    
    // 全局中间件 / Global middleware
    h.Use(LoggingMiddleware())
    h.Use(RecoveryMiddleware())
    h.Use(CORSMiddleware())
    
    // 公开路由 / Public routes
    h.GET("/public", func(ctx context.Context, c *app.RequestContext) {
        c.JSON(consts.StatusOK, map[string]string{"message": "public"})
    })
    
    // 需要认证的路由组 / Routes requiring auth
    authorized := h.Group("/api")
    authorized.Use(AuthMiddleware())
    {
        authorized.GET("/profile", func(ctx context.Context, c *app.RequestContext) {
            userID, _ := c.Get("user_id")
            c.JSON(consts.StatusOK, map[string]interface{}{
                "user_id": userID,
            })
        })
    }
    
    h.Spin()
}
```

## 7. 参数验证 / Parameter Validation

```go
package main

import (
    "context"
    
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/protocol/consts"
)

// 使用 vd 标签进行验证 / Use vd tag for validation
type RegisterReq struct {
    Username string `json:"username" vd:"len($) >= 3 && len($) <= 20; msg:'用户名长度必须在3-20之间'"`
    Password string `json:"password" vd:"len($) >= 6; msg:'密码长度至少6位'"`
    Email    string `json:"email" vd:"email($); msg:'邮箱格式不正确'"`
    Age      int    `json:"age" vd:"$ >= 0 && $ <= 150; msg:'年龄必须在0-150之间'"`
    Phone    string `json:"phone" vd:"regexp('^1[3-9]\\d{9}$'); msg:'手机号格式不正确'"`
}

type UpdateUserReq struct {
    Name  *string `json:"name" vd:"@:len($) >= 2; msg:'名字至少2个字符'"`
    Email *string `json:"email" vd:"@:email($); msg:'邮箱格式不正确'"`
}

func main() {
    h := server.Default()
    
    h.POST("/register", func(ctx context.Context, c *app.RequestContext) {
        var req RegisterReq
        if err := c.BindAndValidate(&req); err != nil {
            c.JSON(consts.StatusBadRequest, map[string]string{
                "error": err.Error(),
            })
            return
        }
        
        c.JSON(consts.StatusOK, map[string]interface{}{
            "message": "注册成功",
            "data":    req,
        })
    })
    
    h.PUT("/user", func(ctx context.Context, c *app.RequestContext) {
        var req UpdateUserReq
        if err := c.BindAndValidate(&req); err != nil {
            c.JSON(consts.StatusBadRequest, map[string]string{
                "error": err.Error(),
            })
            return
        }
        
        c.JSON(consts.StatusOK, map[string]string{
            "message": "更新成功",
        })
    })
    
    h.Spin()
}
```

## 8. 使用 IDL 生成代码 / Generate Code from IDL

```thrift
// idl/api.thrift
namespace go api

struct User {
    1: i64 id,
    2: string name,
    3: string email
}

struct GetUserReq {
    1: i64 id (api.path="id")
}

struct GetUserResp {
    1: User user
}

struct CreateUserReq {
    1: string name (api.body="name")
    2: string email (api.body="email")
}

struct CreateUserResp {
    1: i64 id
}

service UserService {
    GetUserResp GetUser(1: GetUserReq req) (api.get="/api/user/:id")
    CreateUserResp CreateUser(1: CreateUserReq req) (api.post="/api/user")
}
```

生成代码 / Generate code:
```bash
hz new -module your-module-name -idl idl/api.thrift
```

## 9. 配置与优化 / Configuration & Optimization

```go
package main

import (
    "time"
    
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/network/standard"
)

func main() {
    h := server.New(
        // 服务地址 / Server address
        server.WithHostPorts(":8080"),
        
        // 最大请求体大小 / Max request body size
        server.WithMaxRequestBodySize(4 * 1024 * 1024),  // 4MB
        
        // 读写超时 / Read/write timeout
        server.WithReadTimeout(time.Second * 10),
        server.WithWriteTimeout(time.Second * 10),
        server.WithIdleTimeout(time.Minute * 2),
        
        // 使用标准网络库 (跨平台) / Use standard network (cross-platform)
        server.WithTransport(standard.NewTransporter),
        
        // 优雅关闭超时 / Graceful shutdown timeout
        server.WithExitWaitTime(time.Second * 5),
        
        // 禁用默认日期头 / Disable default date header
        server.WithNoDefaultDate(true),
        
        // 启用 H2C (HTTP/2 without TLS)
        server.WithH2C(true),
        
        // TLS 配置 / TLS configuration
        // server.WithTLS(certFile, keyFile),
    )
    
    h.GET("/", func(ctx context.Context, c *app.RequestContext) {
        c.JSON(200, map[string]string{"message": "ok"})
    })
    
    h.Spin()
}
```

## 10. 与其他组件集成 / Integration with Other Components

```go
package main

import (
    "context"
    
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    
    // Prometheus 监控
    "github.com/hertz-contrib/monitor-prometheus"
    
    // OpenTelemetry 链路追踪
    hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
    
    // JWT 认证
    "github.com/hertz-contrib/jwt"
    
    // 限流
    "github.com/hertz-contrib/limiter"
)

// JWT 认证 / JWT authentication
func setupJWT(h *server.Hertz) {
    authMiddleware, _ := jwt.New(&jwt.HertzJWTMiddleware{
        Key:        []byte("secret-key"),
        Timeout:    time.Hour,
        MaxRefresh: time.Hour * 24,
        Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
            var loginReq struct {
                Username string `json:"username"`
                Password string `json:"password"`
            }
            if err := c.BindJSON(&loginReq); err != nil {
                return nil, err
            }
            // 验证用户名密码...
            return loginReq.Username, nil
        },
    })
    
    h.POST("/login", authMiddleware.LoginHandler)
    
    auth := h.Group("/auth")
    auth.Use(authMiddleware.MiddlewareFunc())
    {
        auth.GET("/refresh", authMiddleware.RefreshHandler)
        auth.GET("/protected", func(ctx context.Context, c *app.RequestContext) {
            c.JSON(200, map[string]string{"message": "protected"})
        })
    }
}

// Prometheus 监控 / Prometheus monitoring
func setupPrometheus(h *server.Hertz) {
    h.Use(prometheus.New())
}

// 限流 / Rate limiting
func setupRateLimiter(h *server.Hertz) {
    h.Use(limiter.AdaptiveLimit())
}

func main() {
    h := server.Default()
    
    setupPrometheus(h)
    setupRateLimiter(h)
    setupJWT(h)
    
    h.Spin()
}
```

## 11. 最佳实践 / Best Practices

```go
package main

// 1. 统一响应结构 / Unified response structure
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(c *app.RequestContext, data interface{}) {
    c.JSON(consts.StatusOK, Response{
        Code:    0,
        Message: "success",
        Data:    data,
    })
}

func Error(c *app.RequestContext, code int, message string) {
    c.JSON(consts.StatusOK, Response{
        Code:    code,
        Message: message,
    })
}

// 2. 分层架构 / Layered architecture
// - handler: 处理请求和响应
// - service: 业务逻辑
// - repository: 数据访问

// 3. 使用接口进行依赖注入 / Use interfaces for DI
type UserService interface {
    GetUser(ctx context.Context, id int64) (*User, error)
    CreateUser(ctx context.Context, user *User) error
}

type UserHandler struct {
    userService UserService
}

func NewUserHandler(us UserService) *UserHandler {
    return &UserHandler{userService: us}
}

// 4. 错误处理 / Error handling
type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    return e.Message
}

func HandleError(c *app.RequestContext, err error) {
    if appErr, ok := err.(*AppError); ok {
        Error(c, appErr.Code, appErr.Message)
        return
    }
    Error(c, -1, "Internal Server Error")
}

// 5. 配置管理 / Configuration management
type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        DSN string `yaml:"dsn"`
    } `yaml:"database"`
}
```
