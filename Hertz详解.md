# Hertz 详解与实战

## 目录
- [简介](#简介)
- [Hertz CLI 命令行详解](#hertz-cli-命令行详解)
- [Context 详解](#context-详解)
- [路由与中间件](#路由与中间件)
- [请求与响应](#请求与响应)
- [参数绑定与验证](#参数绑定与验证)
- [渲染](#渲染)
- [客户端](#客户端)
- [中间件开发](#中间件开发)
- [性能优化](#性能优化)
- [实战案例](#实战案例)

---

## 简介

Hertz 是 CloudWeGo 开源的高性能、高可用、易扩展的 Go HTTP 框架。它参考了 Gin、Fasthttp 等优秀框架的设计，在性能、易用性上做了大量优化。

### 主要特性

- **高性能**：基于自研的高性能网络库 Netpoll，相比 Go 标准库 net 性能更优
- **易扩展**：提供丰富的扩展接口，支持自定义序列化、网络库、协议等
- **多协议**：支持 HTTP/1.1、HTTP/2、HTTP/3（QUIC）
- **服务发现**：与服务注册发现组件无缝集成
- **可观测性**：完善的监控、链路追踪支持
- **丰富的中间件**：提供大量开箱即用的中间件

### 安装

```bash
# 安装 Hertz
go get github.com/cloudwego/hertz

# 安装 Hertz 命令行工具
go install github.com/cloudwego/hertz/cmd/hz@latest

# 验证安装
hz --version
```

---

## Hertz CLI 命令行详解

Hertz 提供了 `hz` 命令行工具，用于快速生成项目脚手架和代码。

### hz 命令总览

```bash
hz --help
```

输出：
```
NAME:
   hz - A idl parser and code generator for Hertz projects

USAGE:
   hz [global options] command [command options] [arguments...]

VERSION:
   v0.x.x

COMMANDS:
   new      Generate a new Hertz project
   update   Update an existing Hertz project
   model    Generate model code
   client   Generate client code
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose      turn on verbose mode (default: false)
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

---

### hz new - 创建新项目

#### 基本用法

```bash
# 1. 创建基本项目
hz new -module github.com/myuser/myproject

# 2. 从 Thrift IDL 创建项目
hz new -module github.com/myuser/myproject -idl idl/api.thrift

# 3. 从 Protobuf IDL 创建项目
hz new -module github.com/myuser/myproject -idl idl/api.proto

# 4. 指定服务名
hz new -module github.com/myuser/myproject -service myservice

# 5. 指定输出目录
hz new -module github.com/myuser/myproject -idl idl/api.thrift -out ./output
```

#### 完整参数列表

```bash
hz new [options]

OPTIONS:
   --module value, -m value          指定 Go module 名称（必需）
   --idl value                       指定 IDL 文件路径（thrift 或 proto）
   --service value                   指定服务名称
   --out value, -o value             指定输出目录（默认：当前目录）
   --handler_dir value               指定 handler 目录（默认：biz/handler）
   --model_dir value                 指定 model 目录（默认：biz/model）
   --router_dir value                指定 router 目录（默认：biz/router）
   --client_dir value                指定 client 目录（默认：biz/client）
   --use value                       指定生成的 handler 中的 import 路径
   --proto_path value, -I value      指定 proto 文件的搜索路径（可多次指定）
   --thriftgo value                  指定 thriftgo 的路径
   --protoc value                    指定 protoc 的路径
   --no_recurse                      不递归搜索 include 路径（默认：false）
   --json_enumstr                    使用字符串而非数字表示枚举（默认：false）
   --unset_omitempty                 生成的 struct tag 不包含 omitempty（默认：false）
   --pb_camel_json_tag               生成 protobuf 的 json tag 使用驼峰（默认：false）
   --snake_tag                       生成 snake_case 的 json tag（默认：false）
   --exclude_file value, -E value    指定不生成的文件列表
   --customize_layout value          指定自定义的项目布局文件
   --customize_package value         指定自定义的包名映射文件
```

#### 实战示例

##### 示例 1：创建 RESTful API 项目

**api.thrift**
```thrift
namespace go api

struct User {
    1: i64 id
    2: string username
    3: string email
    4: i32 age
}

struct CreateUserRequest {
    1: string username (api.body="username")
    2: string email (api.body="email")
    3: i32 age (api.body="age")
}

struct CreateUserResponse {
    1: i64 code
    2: string message
    3: User data
}

struct GetUserRequest {
    1: i64 id (api.path="id")
}

struct GetUserResponse {
    1: i64 code
    2: string message
    3: User data
}

struct ListUsersRequest {
    1: i32 page (api.query="page")
    2: i32 page_size (api.query="page_size")
}

struct ListUsersResponse {
    1: i64 code
    2: string message
    3: list<User> data
    4: i64 total
}

service UserService {
    CreateUserResponse CreateUser(1: CreateUserRequest req) (api.post="/api/users")
    GetUserResponse GetUser(1: GetUserRequest req) (api.get="/api/users/:id")
    ListUsersResponse ListUsers(1: ListUsersRequest req) (api.get="/api/users")
}
```

**生成项目**
```bash
hz new -module github.com/myuser/user-service \
    -idl idl/api.thrift \
    -service user-service
```

生成的目录结构：
```
user-service/
├── biz/
│   ├── handler/
│   │   └── user/
│   │       └── user_service.go    # 业务处理器
│   ├── model/
│   │   └── api/
│   │       └── api.go             # 数据模型
│   └── router/
│       ├── register.go            # 路由注册
│       └── user/
│           └── user.go            # 用户路由
├── go.mod
├── go.sum
├── main.go                        # 主入口
├── router.go                      # 路由配置
└── router_gen.go                  # 生成的路由代码
```

##### 示例 2：使用 Protobuf IDL

**api.proto**
```protobuf
syntax = "proto3";

package api;

option go_package = "github.com/myuser/myproject/biz/model/api";

import "api/annotations.proto";

message User {
  int64 id = 1;
  string username = 2;
  string email = 3;
  int32 age = 4;
}

message CreateUserRequest {
  string username = 1;
  string email = 2;
  int32 age = 3;
}

message CreateUserResponse {
  int64 code = 1;
  string message = 2;
  User data = 3;
}

service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {
    option (api.post) = "/api/users";
  }
}
```

**生成项目**
```bash
hz new -module github.com/myuser/user-service \
    -idl idl/api.proto \
    -proto_path idl \
    -service user-service
```

---

### hz update - 更新项目

当 IDL 文件变更后，使用 `hz update` 更新代码。

#### 基本用法

```bash
# 1. 更新项目（在项目根目录执行）
hz update -idl idl/api.thrift

# 2. 更新并重新生成所有文件
hz update -idl idl/api.thrift --force

# 3. 只更新 handler
hz update -idl idl/api.thrift --handler_dir biz/handler
```

#### 完整参数

```bash
hz update [options]

OPTIONS:
   --idl value                       指定 IDL 文件路径
   --handler_dir value               指定 handler 目录
   --model_dir value                 指定 model 目录
   --router_dir value                指定 router 目录
   --force                           强制更新，覆盖已修改的文件（默认：false）
   其他参数与 hz new 相同
```

---

### hz model - 生成模型代码

单独生成数据模型代码。

```bash
# 从 thrift 生成模型
hz model -idl idl/api.thrift -model_dir biz/model

# 从 proto 生成模型
hz model -idl idl/api.proto -model_dir biz/model -proto_path idl
```

---

### 常用命令组合

#### 1. 快速创建微服务项目

```bash
#!/bin/bash

# 项目配置
PROJECT_NAME="user-service"
MODULE_NAME="github.com/mycompany/${PROJECT_NAME}"
SERVICE_NAME="${PROJECT_NAME}"

# 创建项目目录
mkdir -p ${PROJECT_NAME}
cd ${PROJECT_NAME}

# 创建 IDL 目录
mkdir -p idl

# 创建 IDL 文件
cat > idl/api.thrift << 'EOF'
namespace go api

// 用户模型
struct User {
    1: i64 id
    2: string username
    3: string email
    4: i32 age
    5: string created_at
}

// 创建用户请求
struct CreateUserRequest {
    1: string username (api.body="username")
    2: string email (api.body="email")
    3: i32 age (api.body="age")
}

// 通用响应
struct Response {
    1: i64 code
    2: string message
    3: optional User data
    4: optional list<User> list
    5: optional i64 total
}

// 用户服务
service UserService {
    Response CreateUser(1: CreateUserRequest req) (api.post="/api/v1/users")
    Response GetUser(1: i64 id) (api.get="/api/v1/users/:id")
    Response UpdateUser(1: i64 id, 2: CreateUserRequest req) (api.put="/api/v1/users/:id")
    Response DeleteUser(1: i64 id) (api.delete="/api/v1/users/:id")
    Response ListUsers(1: i32 page, 2: i32 page_size) (api.get="/api/v1/users")
}
EOF

# 生成项目代码
hz new -module ${MODULE_NAME} \
    -idl idl/api.thrift \
    -service ${SERVICE_NAME}

# 下载依赖
go mod tidy

echo "项目创建成功！"
echo "目录: $(pwd)"
echo "启动命令: go run ."
```

#### 2. 多服务项目生成

```bash
#!/bin/bash

# 创建多个服务
services=("user-service" "post-service" "comment-service")

for service in "${services[@]}"; do
    hz new -module github.com/mycompany/${service} \
        -idl idl/${service}.thrift \
        -service ${service} \
        -out ./${service}
done
```

---

## Context 详解

Hertz 的 `RequestContext` 是处理请求的核心，类似于 Gin 的 Context。

### Context 结构

```go
type RequestContext struct {
    // 包含请求和响应的所有信息
}
```

### 获取请求信息

#### 1. 获取请求参数

```go
package handler

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/common/utils"
)

// 路径参数
func GetPathParam(ctx context.Context, c *app.RequestContext) {
    // GET /users/:id
    id := c.Param("id")
    c.JSON(200, utils.H{"id": id})
}

// Query 参数
func GetQueryParam(ctx context.Context, c *app.RequestContext) {
    // GET /users?name=zhangsan&age=25
    name := c.Query("name")           // 获取单个参数
    age := c.DefaultQuery("age", "0") // 带默认值
    
    // 获取所有 query 参数
    query := c.Request.URI().QueryArgs()
    
    c.JSON(200, utils.H{
        "name": name,
        "age":  age,
    })
}

// Form 参数（application/x-www-form-urlencoded）
func GetFormParam(ctx context.Context, c *app.RequestContext) {
    // POST /login
    username := c.PostForm("username")
    password := c.DefaultPostForm("password", "")
    
    c.JSON(200, utils.H{
        "username": username,
    })
}

// 获取 Body（JSON）
func GetBody(ctx context.Context, c *app.RequestContext) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    
    if err := c.Bind(&req); err != nil {
        c.JSON(400, utils.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, utils.H{"received": req})
}

// 获取 Header
func GetHeader(ctx context.Context, c *app.RequestContext) {
    // 获取单个 header
    auth := c.GetHeader("Authorization")
    
    // 获取所有 headers
    headers := &c.Request.Header
    
    c.JSON(200, utils.H{
        "auth": string(auth),
    })
}

// 获取 Cookie
func GetCookie(ctx context.Context, c *app.RequestContext) {
    cookie := c.Cookie("session_id")
    c.JSON(200, utils.H{"session_id": string(cookie)})
}
```

#### 2. 文件上传

```go
// 单文件上传
func UploadFile(ctx context.Context, c *app.RequestContext) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, utils.H{"error": err.Error()})
        return
    }
    
    // 保存文件
    dst := "/path/to/save/" + file.Filename
    if err := c.SaveUploadedFile(file, dst); err != nil {
        c.JSON(500, utils.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, utils.H{
        "filename": file.Filename,
        "size":     file.Size,
    })
}

// 多文件上传
func UploadMultipleFiles(ctx context.Context, c *app.RequestContext) {
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(400, utils.H{"error": err.Error()})
        return
    }
    
    files := form.File["files"]
    
    for _, file := range files {
        dst := "/path/to/save/" + file.Filename
        if err := c.SaveUploadedFile(file, dst); err != nil {
            c.JSON(500, utils.H{"error": err.Error()})
            return
        }
    }
    
    c.JSON(200, utils.H{
        "count": len(files),
    })
}
```

### 设置响应

#### 1. 响应 JSON

```go
func ResponseJSON(ctx context.Context, c *app.RequestContext) {
    // 方式 1：使用 utils.H
    c.JSON(200, utils.H{
        "code":    0,
        "message": "success",
        "data":    "hello",
    })
    
    // 方式 2：使用 struct
    type Response struct {
        Code    int         `json:"code"`
        Message string      `json:"message"`
        Data    interface{} `json:"data"`
    }
    
    c.JSON(200, Response{
        Code:    0,
        Message: "success",
        Data:    "hello",
    })
}
```

#### 2. 响应其他格式

```go
// XML
func ResponseXML(ctx context.Context, c *app.RequestContext) {
    type User struct {
        Name string `xml:"name"`
        Age  int    `xml:"age"`
    }
    
    c.XML(200, User{Name: "zhangsan", Age: 25})
}

// HTML
func ResponseHTML(ctx context.Context, c *app.RequestContext) {
    c.HTML(200, "<h1>Hello World</h1>")
}

// String
func ResponseString(ctx context.Context, c *app.RequestContext) {
    c.String(200, "Hello %s", "World")
}

// Data（二进制）
func ResponseData(ctx context.Context, c *app.RequestContext) {
    data := []byte("binary data")
    c.Data(200, "application/octet-stream", data)
}

// File
func ResponseFile(ctx context.Context, c *app.RequestContext) {
    c.File("/path/to/file.pdf")
}

// FileAttachment（下载）
func DownloadFile(ctx context.Context, c *app.RequestContext) {
    c.FileAttachment("/path/to/file.pdf", "download.pdf")
}

// Redirect
func RedirectURL(ctx context.Context, c *app.RequestContext) {
    c.Redirect(302, []byte("https://example.com"))
}
```

#### 3. 设置 Header 和 Cookie

```go
func SetHeaderAndCookie(ctx context.Context, c *app.RequestContext) {
    // 设置 Header
    c.Header("X-Custom-Header", "value")
    c.Response.Header.Set("Content-Type", "application/json")
    
    // 设置 Cookie
    c.SetCookie(
        "session_id",           // name
        "abc123",               // value
        3600,                   // maxAge (seconds)
        "/",                    // path
        "example.com",          // domain
        protocol.CookieSameSiteDefaultMode, // sameSite
        true,                   // secure
        true,                   // httpOnly
    )
    
    c.JSON(200, utils.H{"message": "headers and cookies set"})
}
```

### Context 存储

#### 使用 Set/Get 存储数据

```go
// 中间件中设置数据
func AuthMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 验证 token
        token := c.GetHeader("Authorization")
        
        // 假设验证成功，存储用户信息
        userID := 123
        c.Set("user_id", userID)
        c.Set("username", "zhangsan")
        
        c.Next(ctx)
    }
}

// Handler 中获取数据
func GetUserInfo(ctx context.Context, c *app.RequestContext) {
    // 获取存储的数据
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(401, utils.H{"error": "unauthorized"})
        return
    }
    
    // 类型断言
    uid := userID.(int)
    
    // 或使用 MustGet（不存在会 panic）
    username := c.MustGet("username").(string)
    
    c.JSON(200, utils.H{
        "user_id":  uid,
        "username": username,
    })
}
```

### Context 其他方法

```go
func ContextMethods(ctx context.Context, c *app.RequestContext) {
    // 获取请求方法
    method := c.Method()
    
    // 获取请求路径
    path := string(c.Path())
    
    // 获取完整 URL
    url := c.URI().String()
    
    // 获取 Host
    host := string(c.Host())
    
    // 获取客户端 IP
    ip := c.ClientIP()
    
    // 获取 User-Agent
    ua := string(c.UserAgent())
    
    // 获取 Content-Type
    ct := string(c.ContentType())
    
    // 判断是否是 Websocket 请求
    isWs := c.IsGet() && c.GetHeader("Upgrade") == []byte("websocket")
    
    // 中止请求处理
    c.Abort()
    
    // 中止并返回错误
    c.AbortWithMsg("error message", 500)
    
    // 中止并返回 JSON
    c.AbortWithStatusJSON(400, utils.H{"error": "bad request"})
    
    c.JSON(200, utils.H{
        "method": method,
        "path":   path,
        "url":    url,
        "host":   host,
        "ip":     ip,
        "ua":     ua,
        "ct":     ct,
    })
}
```

---

## 路由与中间件

### 基本路由

```go
package main

import (
    "context"
    
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/common/utils"
)

func main() {
    h := server.Default()
    
    // GET 请求
    h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
        c.JSON(200, utils.H{"message": "pong"})
    })
    
    // POST 请求
    h.POST("/users", createUser)
    
    // PUT 请求
    h.PUT("/users/:id", updateUser)
    
    // DELETE 请求
    h.DELETE("/users/:id", deleteUser)
    
    // PATCH 请求
    h.PATCH("/users/:id", patchUser)
    
    // HEAD 请求
    h.HEAD("/users/:id", headUser)
    
    // OPTIONS 请求
    h.OPTIONS("/users/:id", optionsUser)
    
    // 任意方法
    h.Any("/any", anyHandler)
    
    // 静态文件
    h.Static("/static", "./static")
    h.StaticFile("/favicon.ico", "./favicon.ico")
    
    h.Spin()
}

func createUser(ctx context.Context, c *app.RequestContext) {
    c.JSON(201, utils.H{"message": "user created"})
}

func updateUser(ctx context.Context, c *app.RequestContext) {
    id := c.Param("id")
    c.JSON(200, utils.H{"message": "user updated", "id": id})
}

func deleteUser(ctx context.Context, c *app.RequestContext) {
    id := c.Param("id")
    c.JSON(200, utils.H{"message": "user deleted", "id": id})
}

func patchUser(ctx context.Context, c *app.RequestContext) {
    c.JSON(200, utils.H{"message": "user patched"})
}

func headUser(ctx context.Context, c *app.RequestContext) {
    c.Status(200)
}

func optionsUser(ctx context.Context, c *app.RequestContext) {
    c.Status(200)
}

func anyHandler(ctx context.Context, c *app.RequestContext) {
    c.JSON(200, utils.H{"method": string(c.Method())})
}
```

### 路由分组

```go
func SetupRoutes(h *server.Hertz) {
    // API v1 分组
    v1 := h.Group("/api/v1")
    {
        // 用户路由
        users := v1.Group("/users")
        {
            users.GET("", listUsers)
            users.GET("/:id", getUser)
            users.POST("", createUser)
            users.PUT("/:id", updateUser)
            users.DELETE("/:id", deleteUser)
        }
        
        // 文章路由
        posts := v1.Group("/posts")
        {
            posts.GET("", listPosts)
            posts.GET("/:id", getPost)
            posts.POST("", createPost)
            posts.PUT("/:id", updatePost)
            posts.DELETE("/:id", deletePost)
        }
    }
    
    // API v2 分组
    v2 := h.Group("/api/v2")
    {
        v2.GET("/users", listUsersV2)
    }
}

func listUsers(ctx context.Context, c *app.RequestContext) {
    c.JSON(200, utils.H{"users": []string{}})
}

func getUser(ctx context.Context, c *app.RequestContext) {
    c.JSON(200, utils.H{"user": utils.H{}})
}

func listPosts(ctx context.Context, c *app.RequestContext) {
    c.JSON(200, utils.H{"posts": []string{}})
}

func getPost(ctx context.Context, c *app.RequestContext) {
    c.JSON(200, utils.H{"post": utils.H{}})
}

func createPost(ctx context.Context, c *app.RequestContext) {
    c.JSON(201, utils.H{"message": "created"})
}

func updatePost(ctx context.Context, c *app.RequestContext) {
    c.JSON(200, utils.H{"message": "updated"})
}

func deletePost(ctx context.Context, c *app.RequestContext) {
    c.JSON(200, utils.H{"message": "deleted"})
}

func listUsersV2(ctx context.Context, c *app.RequestContext) {
    c.JSON(200, utils.H{"users": []string{}, "version": "v2"})
}
```

### 中间件

#### 全局中间件

```go
func main() {
    h := server.Default() // 包含 Logger 和 Recovery 中间件
    
    // 或者创建空的服务器
    h := server.New()
    
    // 添加全局中间件
    h.Use(Logger())
    h.Use(Recovery())
    h.Use(CORS())
    
    h.Spin()
}
```

#### 路由级中间件

```go
func main() {
    h := server.Default()
    
    // 单个路由使用中间件
    h.GET("/admin", AuthMiddleware(), adminHandler)
    
    // 分组使用中间件
    api := h.Group("/api", RateLimitMiddleware())
    {
        api.GET("/users", listUsers)
    }
    
    // 多个中间件
    auth := h.Group("/auth")
    auth.Use(AuthMiddleware(), LoggerMiddleware())
    {
        auth.GET("/profile", getProfile)
    }
    
    h.Spin()
}

func adminHandler(ctx context.Context, c *app.RequestContext) {
    c.JSON(200, utils.H{"message": "admin page"})
}

func getProfile(ctx context.Context, c *app.RequestContext) {
    c.JSON(200, utils.H{"profile": utils.H{}})
}
```

#### 自定义中间件

```go
// 日志中间件
func LoggerMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        start := time.Now()
        
        // 处理请求
        c.Next(ctx)
        
        // 请求后
        latency := time.Since(start)
        statusCode := c.Response.StatusCode()
        
        log.Printf("[%d] %s %s %v",
            statusCode,
            c.Method(),
            c.Path(),
            latency,
        )
    }
}

// 认证中间件
func AuthMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        token := string(c.GetHeader("Authorization"))
        
        if token == "" {
            c.AbortWithStatusJSON(401, utils.H{
                "error": "missing authorization header",
            })
            return
        }
        
        // 验证 token
        userID, err := validateToken(token)
        if err != nil {
            c.AbortWithStatusJSON(401, utils.H{
                "error": "invalid token",
            })
            return
        }
        
        // 存储用户信息
        c.Set("user_id", userID)
        
        c.Next(ctx)
    }
}

// CORS 中间件
func CORS() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if string(c.Method()) == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next(ctx)
    }
}

// 限流中间件
func RateLimitMiddleware() app.HandlerFunc {
    limiter := rate.NewLimiter(10, 100) // 每秒 10 个请求，桶容量 100
    
    return func(ctx context.Context, c *app.RequestContext) {
        if !limiter.Allow() {
            c.AbortWithStatusJSON(429, utils.H{
                "error": "too many requests",
            })
            return
        }
        
        c.Next(ctx)
    }
}

// Recovery 中间件
func Recovery() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic recovered: %v", err)
                c.AbortWithStatusJSON(500, utils.H{
                    "error": "internal server error",
                })
            }
        }()
        
        c.Next(ctx)
    }
}

func validateToken(token string) (int, error) {
    // 实现 token 验证逻辑
    if token == "valid_token" {
        return 123, nil
    }
    return 0, errors.New("invalid token")
}
```

---

## 参数绑定与验证

### 参数绑定

```go
package handler

import (
    "context"
    
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/common/utils"
)

// 绑定 JSON
type CreateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Age      int    `json:"age" binding:"required,gte=0,lte=150"`
}

func CreateUser(ctx context.Context, c *app.RequestContext) {
    var req CreateUserRequest
    
    // 绑定并验证
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, utils.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, utils.H{
        "message": "user created",
        "data":    req,
    })
}

// 绑定 Query 参数
type ListUsersQuery struct {
    Page     int    `query:"page" binding:"required,gte=1"`
    PageSize int    `query:"page_size" binding:"required,gte=1,lte=100"`
    Keyword  string `query:"keyword"`
}

func ListUsers(ctx context.Context, c *app.RequestContext) {
    var query ListUsersQuery
    
    if err := c.BindAndValidate(&query); err != nil {
        c.JSON(400, utils.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, utils.H{
        "page":      query.Page,
        "page_size": query.PageSize,
        "keyword":   query.Keyword,
    })
}

// 绑定路径参数
type GetUserURI struct {
    ID int `path:"id" binding:"required,gte=1"`
}

func GetUser(ctx context.Context, c *app.RequestContext) {
    var uri GetUserURI
    
    if err := c.BindAndValidate(&uri); err != nil {
        c.JSON(400, utils.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, utils.H{
        "user_id": uri.ID,
    })
}

// 绑定 Form 数据
type LoginForm struct {
    Username string `form:"username" binding:"required"`
    Password string `form:"password" binding:"required,min=6"`
}

func Login(ctx context.Context, c *app.RequestContext) {
    var form LoginForm
    
    if err := c.BindAndValidate(&form); err != nil {
        c.JSON(400, utils.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, utils.H{
        "message": "login successful",
    })
}

// 绑定文件上传
type UploadRequest struct {
    File     *multipart.FileHeader `form:"file" binding:"required"`
    Category string                `form:"category" binding:"required"`
}

func UploadFile(ctx context.Context, c *app.RequestContext) {
    var req UploadRequest
    
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, utils.H{"error": err.Error()})
        return
    }
    
    // 保存文件
    dst := "/uploads/" + req.File.Filename
    if err := c.SaveUploadedFile(req.File, dst); err != nil {
        c.JSON(500, utils.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, utils.H{
        "filename": req.File.Filename,
        "category": req.Category,
    })
}
```

### 验证标签

```go
type User struct {
    // 必填
    Username string `json:"username" binding:"required"`
    
    // 邮箱格式
    Email string `json:"email" binding:"required,email"`
    
    // 数值范围
    Age int `json:"age" binding:"gte=0,lte=150"`
    
    // 字符串长度
    Password string `json:"password" binding:"min=6,max=20"`
    
    // 枚举值
    Gender string `json:"gender" binding:"oneof=male female other"`
    
    // URL 格式
    Website string `json:"website" binding:"url"`
    
    // IP 地址
    IPAddress string `json:"ip" binding:"ip"`
    
    // 正则表达式
    Phone string `json:"phone" binding:"regexp=^1[3-9]\\d{9}$"`
    
    // 自定义验证
    Code string `json:"code" binding:"len=6,numeric"`
    
    // 嵌套验证
    Profile Profile `json:"profile" binding:"required"`
}

type Profile struct {
    Bio      string `json:"bio" binding:"max=500"`
    Location string `json:"location"`
}
```

### 自定义验证器

```go
import (
    "github.com/cloudwego/hertz/pkg/app/server/binding"
    "github.com/go-playground/validator/v10"
)

// 自定义验证函数
func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    
    // 用户名只能包含字母、数字和下划线
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
    return matched
}

func init() {
    // 注册自定义验证器
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("username", validateUsername)
    }
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,username"`
    Email    string `json:"email" binding:"required,email"`
}
```

---

## 渲染

### HTML 渲染

```go
package main

import (
    "context"
    "html/template"
    
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
    h := server.Default()
    
    // 加载 HTML 模板
    h.LoadHTMLGlob("templates/*")
    
    // 或加载指定文件
    // h.LoadHTMLFiles("templates/index.html", "templates/user.html")
    
    h.GET("/", indexHandler)
    h.GET("/user/:id", userHandler)
    
    h.Spin()
}

func indexHandler(ctx context.Context, c *app.RequestContext) {
    c.HTML(200, "index.html", utils.H{
        "title": "Home Page",
        "user":  "zhangsan",
    })
}

func userHandler(ctx context.Context, c *app.RequestContext) {
    id := c.Param("id")
    
    c.HTML(200, "user.html", utils.H{
        "title":   "User Profile",
        "user_id": id,
        "user": map[string]interface{}{
            "name":  "zhangsan",
            "email": "zhangsan@example.com",
        },
    })
}
```

**templates/index.html**
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .title }}</title>
</head>
<body>
    <h1>Welcome, {{ .user }}!</h1>
</body>
</html>
```

**templates/user.html**
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .title }}</title>
</head>
<body>
    <h1>User Profile</h1>
    <p>ID: {{ .user_id }}</p>
    <p>Name: {{ .user.name }}</p>
    <p>Email: {{ .user.email }}</p>
</body>
</html>
```

### 自定义模板函数

```go
func main() {
    h := server.Default()
    
    // 设置自定义函数
    h.SetFuncMap(template.FuncMap{
        "upper": strings.ToUpper,
        "lower": strings.ToLower,
        "formatDate": func(t time.Time) string {
            return t.Format("2006-01-02 15:04:05")
        },
    })
    
    h.LoadHTMLGlob("templates/*")
    
    h.GET("/", func(ctx context.Context, c *app.RequestContext) {
        c.HTML(200, "index.html", utils.H{
            "title": "hello world",
            "time":  time.Now(),
        })
    })
    
    h.Spin()
}
```

**templates/index.html**
```html
<h1>{{ upper .title }}</h1>
<p>Time: {{ formatDate .time }}</p>
```

---

## 客户端

Hertz 提供了高性能的 HTTP 客户端。

### 基本使用

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/cloudwego/hertz/pkg/app/client"
    "github.com/cloudwego/hertz/pkg/protocol"
)

func main() {
    c, err := client.NewClient()
    if err != nil {
        panic(err)
    }
    
    ctx := context.Background()
    
    // GET 请求
    status, body, err := c.Get(ctx, nil, "http://example.com/api/users")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Status: %d, Body: %s\n", status, body)
    
    // POST 请求
    reqBody := []byte(`{"username":"zhangsan","email":"zhangsan@example.com"}`)
    status, body, err = c.Post(ctx, nil, "http://example.com/api/users", reqBody)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Status: %d, Body: %s\n", status, body)
}
```

### 高级用法

```go
package main

import (
    "context"
    "time"
    
    "github.com/cloudwego/hertz/pkg/app/client"
    "github.com/cloudwego/hertz/pkg/protocol"
    "github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
    // 创建客户端配置
    c, err := client.NewClient(
        client.WithMaxConnsPerHost(100),
        client.WithClientReadTimeout(10*time.Second),
        client.WithDialTimeout(5*time.Second),
    )
    if err != nil {
        panic(err)
    }
    
    ctx := context.Background()
    
    // 创建请求
    req := &protocol.Request{}
    res := &protocol.Response{}
    
    // 设置请求
    req.SetMethod(consts.MethodPost)
    req.SetRequestURI("http://example.com/api/users")
    
    // 设置 Header
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer token")
    
    // 设置 Body
    req.SetBody([]byte(`{"username":"zhangsan"}`))
    
    // 发送请求
    err = c.Do(ctx, req, res)
    if err != nil {
        panic(err)
    }
    
    // 处理响应
    fmt.Printf("Status: %d\n", res.StatusCode())
    fmt.Printf("Body: %s\n", res.Body())
}
```

### 文件上传

```go
func UploadFile(ctx context.Context) error {
    c, _ := client.NewClient()
    
    req := &protocol.Request{}
    res := &protocol.Response{}
    
    req.SetMethod(consts.MethodPost)
    req.SetRequestURI("http://example.com/upload")
    
    // 设置 multipart
    req.SetMultipartFormData(map[string]string{
        "field1": "value1",
    })
    
    // 添加文件
    req.SetFile("file", "/path/to/file.txt")
    
    return c.Do(ctx, req, res)
}
```

---

## 中间件开发

### 标准中间件模板

```go
func MiddlewareTemplate() app.HandlerFunc {
    // 初始化代码（只执行一次）
    setup := doSomeSetup()
    
    return func(ctx context.Context, c *app.RequestContext) {
        // 请求前处理
        beforeRequest()
        
        // 调用下一个处理器
        c.Next(ctx)
        
        // 请求后处理
        afterRequest()
    }
}

func doSomeSetup() interface{} {
    // 初始化逻辑
    return nil
}

func beforeRequest() {
    // 请求前逻辑
}

func afterRequest() {
    // 请求后逻辑
}
```

### 实战中间件示例

#### 1. 请求 ID 中间件

```go
func RequestIDMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        requestID := c.GetHeader("X-Request-ID")
        if len(requestID) == 0 {
            requestID = []byte(uuid.New().String())
        }
        
        c.Set("request_id", string(requestID))
        c.Header("X-Request-ID", string(requestID))
        
        c.Next(ctx)
    }
}
```

#### 2. 超时控制中间件

```go
func TimeoutMiddleware(timeout time.Duration) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        ctx, cancel := context.WithTimeout(ctx, timeout)
        defer cancel()
        
        done := make(chan struct{})
        
        go func() {
            c.Next(ctx)
            close(done)
        }()
        
        select {
        case <-done:
            return
        case <-ctx.Done():
            c.AbortWithStatusJSON(504, utils.H{
                "error": "request timeout",
            })
        }
    }
}
```

#### 3. 缓存中间件

```go
func CacheMiddleware(cache *redis.Client, ttl time.Duration) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 只缓存 GET 请求
        if string(c.Method()) != "GET" {
            c.Next(ctx)
            return
        }
        
        // 生成缓存键
        cacheKey := generateCacheKey(c)
        
        // 尝试从缓存获取
        cached, err := cache.Get(ctx, cacheKey).Bytes()
        if err == nil {
            c.Data(200, "application/json", cached)
            c.Abort()
            return
        }
        
        // 捕获响应
        writer := &responseWriter{ResponseWriter: c.Response.BodyWriter()}
        c.Response.SetBodyWriter(writer)
        
        c.Next(ctx)
        
        // 缓存响应
        if c.Response.StatusCode() == 200 {
            cache.Set(ctx, cacheKey, writer.body.Bytes(), ttl)
        }
    }
}

type responseWriter struct {
    http.ResponseWriter
    body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
    w.body.Write(b)
    return w.ResponseWriter.Write(b)
}

func generateCacheKey(c *app.RequestContext) string {
    return fmt.Sprintf("cache:%s:%s", c.Path(), c.Request.URI().QueryString())
}
```

---

## 性能优化

### 1. 使用 Netpoll

```go
package main

import (
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/network/netpoll"
)

func main() {
    // 使用 netpoll 网络库（默认）
    h := server.New(
        server.WithHostPorts(":8080"),
        server.WithTransport(netpoll.NewTransporter),
    )
    
    h.Spin()
}
```

### 2. 连接池配置

```go
func main() {
    h := server.New(
        server.WithHostPorts(":8080"),
        server.WithMaxRequestBodySize(4*1024*1024), // 4MB
        server.WithIdleTimeout(60*time.Second),
        server.WithReadTimeout(3*time.Second),
        server.WithWriteTimeout(3*time.Second),
    )
    
    h.Spin()
}
```

### 3. 零拷贝

```go
// 使用 Nocopy API
func NoCopyHandler(ctx context.Context, c *app.RequestContext) {
    // 获取 body（零拷贝）
    body := c.Request.Body()
    
    // 直接使用，不要修改
    processBody(body)
    
    c.Status(200)
}

func processBody(body []byte) {
    // 处理逻辑
}
```

### 4. 对象池

```go
var requestPool = sync.Pool{
    New: func() interface{} {
        return &Request{}
    },
}

func Handler(ctx context.Context, c *app.RequestContext) {
    req := requestPool.Get().(*Request)
    defer requestPool.Put(req)
    
    // 使用 req
    req.Reset()
    req.Parse(c)
    
    c.JSON(200, utils.H{"message": "ok"})
}

type Request struct {
    ID   int
    Name string
}

func (r *Request) Reset() {
    r.ID = 0
    r.Name = ""
}

func (r *Request) Parse(c *app.RequestContext) {
    c.Bind(r)
}
```

---

## 实战案例

### 完整的 RESTful API 服务

```go
package main

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"
    
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/common/utils"
)

// 用户模型
type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// 内存存储（实际应用中应使用数据库）
type UserStore struct {
    mu    sync.RWMutex
    users map[int]*User
    nextID int
}

func NewUserStore() *UserStore {
    return &UserStore{
        users: make(map[int]*User),
        nextID: 1,
    }
}

func (s *UserStore) Create(user *User) *User {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    user.ID = s.nextID
    user.CreatedAt = time.Now()
    s.users[user.ID] = user
    s.nextID++
    
    return user
}

func (s *UserStore) Get(id int) (*User, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    user, ok := s.users[id]
    return user, ok
}

func (s *UserStore) List() []*User {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    users := make([]*User, 0, len(s.users))
    for _, user := range s.users {
        users = append(users, user)
    }
    
    return users
}

func (s *UserStore) Update(id int, updates *User) (*User, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    user, ok := s.users[id]
    if !ok {
        return nil, false
    }
    
    if updates.Username != "" {
        user.Username = updates.Username
    }
    if updates.Email != "" {
        user.Email = updates.Email
    }
    
    return user, true
}

func (s *UserStore) Delete(id int) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    _, ok := s.users[id]
    if ok {
        delete(s.users, id)
    }
    
    return ok
}

// Handler
type UserHandler struct {
    store *UserStore
}

func NewUserHandler(store *UserStore) *UserHandler {
    return &UserHandler{store: store}
}

// 创建用户
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Email    string `json:"email" binding:"required,email"`
}

func (h *UserHandler) Create(ctx context.Context, c *app.RequestContext) {
    var req CreateUserRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, utils.H{
            "code":    400,
            "message": err.Error(),
        })
        return
    }
    
    user := &User{
        Username: req.Username,
        Email:    req.Email,
    }
    
    user = h.store.Create(user)
    
    c.JSON(201, utils.H{
        "code":    0,
        "message": "success",
        "data":    user,
    })
}

// 获取用户
func (h *UserHandler) Get(ctx context.Context, c *app.RequestContext) {
    id := c.Param("id")
    
    var idInt int
    if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
        c.JSON(400, utils.H{
            "code":    400,
            "message": "invalid user id",
        })
        return
    }
    
    user, ok := h.store.Get(idInt)
    if !ok {
        c.JSON(404, utils.H{
            "code":    404,
            "message": "user not found",
        })
        return
    }
    
    c.JSON(200, utils.H{
        "code":    0,
        "message": "success",
        "data":    user,
    })
}

// 列出用户
func (h *UserHandler) List(ctx context.Context, c *app.RequestContext) {
    users := h.store.List()
    
    c.JSON(200, utils.H{
        "code":    0,
        "message": "success",
        "data":    users,
        "total":   len(users),
    })
}

// 更新用户
type UpdateUserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
}

func (h *UserHandler) Update(ctx context.Context, c *app.RequestContext) {
    id := c.Param("id")
    
    var idInt int
    if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
        c.JSON(400, utils.H{
            "code":    400,
            "message": "invalid user id",
        })
        return
    }
    
    var req UpdateUserRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, utils.H{
            "code":    400,
            "message": err.Error(),
        })
        return
    }
    
    user, ok := h.store.Update(idInt, &User{
        Username: req.Username,
        Email:    req.Email,
    })
    
    if !ok {
        c.JSON(404, utils.H{
            "code":    404,
            "message": "user not found",
        })
        return
    }
    
    c.JSON(200, utils.H{
        "code":    0,
        "message": "success",
        "data":    user,
    })
}

// 删除用户
func (h *UserHandler) Delete(ctx context.Context, c *app.RequestContext) {
    id := c.Param("id")
    
    var idInt int
    if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
        c.JSON(400, utils.H{
            "code":    400,
            "message": "invalid user id",
        })
        return
    }
    
    ok := h.store.Delete(idInt)
    if !ok {
        c.JSON(404, utils.H{
            "code":    404,
            "message": "user not found",
        })
        return
    }
    
    c.JSON(200, utils.H{
        "code":    0,
        "message": "user deleted",
    })
}

// 主函数
func main() {
    // 创建服务器
    h := server.Default(
        server.WithHostPorts(":8080"),
    )
    
    // 创建存储和 handler
    store := NewUserStore()
    userHandler := NewUserHandler(store)
    
    // 注册路由
    api := h.Group("/api/v1")
    {
        users := api.Group("/users")
        {
            users.POST("", userHandler.Create)
            users.GET("", userHandler.List)
            users.GET("/:id", userHandler.Get)
            users.PUT("/:id", userHandler.Update)
            users.DELETE("/:id", userHandler.Delete)
        }
    }
    
    // 健康检查
    h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
        c.JSON(200, utils.H{
            "status": "ok",
            "time":   time.Now().Unix(),
        })
    })
    
    log.Printf("Server started on :8080")
    h.Spin()
}
```

### 测试 API

```bash
# 创建用户
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"zhangsan","email":"zhangsan@example.com"}'

# 获取用户列表
curl http://localhost:8080/api/v1/users

# 获取单个用户
curl http://localhost:8080/api/v1/users/1

# 更新用户
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"username":"zhangsan_updated","email":"new@example.com"}'

# 删除用户
curl -X DELETE http://localhost:8080/api/v1/users/1

# 健康检查
curl http://localhost:8080/health
```

---

## 总结

本文详细介绍了 Hertz 框架的各个方面：

1. **CLI 工具**：hz 命令行工具的完整使用方法
2. **Context**：RequestContext 的详细 API 和使用方法
3. **路由**：路由配置、分组、参数获取
4. **中间件**：内置和自定义中间件开发
5. **参数处理**：请求绑定、验证、响应渲染
6. **客户端**：HTTP 客户端的使用
7. **性能优化**：各种性能优化技巧
8. **实战**：完整的 RESTful API 示例

### 最佳实践

1. 使用 hz 工具生成项目脚手架
2. 合理使用中间件处理通用逻辑
3. 使用参数绑定和验证简化代码
4. 利用 Netpoll 提升性能
5. 使用对象池减少 GC 压力
6. 合理配置超时和限流
7. 完善的错误处理和日志记录
8. 使用服务注册发现实现微服务
9. 集成监控和链路追踪
10. 编写单元测试和集成测试

更多详细信息请参考 [Hertz 官方文档](https://www.cloudwego.io/zh/docs/hertz/)。
