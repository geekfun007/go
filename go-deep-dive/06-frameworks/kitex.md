# Kitex - 字节跳动高性能 RPC 框架 / ByteDance High-Performance RPC Framework

## 1. 简介与安装 / Introduction & Installation

Kitex 是字节跳动开源的高性能、强可扩展的 Go RPC 框架，支持多协议并且拥有丰富的特性。

```bash
# 安装 Kitex 命令行工具 / Install Kitex CLI
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest

# 安装 thriftgo (用于生成代码) / Install thriftgo
go install github.com/cloudwego/thriftgo@latest

# 添加依赖 / Add dependencies
go get github.com/cloudwego/kitex
```

## 2. 定义 IDL / Define IDL

```thrift
// idl/user.thrift
namespace go api.user

struct User {
    1: i64 id,
    2: string name,
    3: string email,
    4: i32 age
}

struct GetUserRequest {
    1: i64 id
}

struct GetUserResponse {
    1: User user
}

struct CreateUserRequest {
    1: string name,
    2: string email,
    3: i32 age
}

struct CreateUserResponse {
    1: i64 id,
    2: bool success
}

struct ListUsersRequest {
    1: i32 page,
    2: i32 page_size
}

struct ListUsersResponse {
    1: list<User> users,
    2: i32 total
}

service UserService {
    GetUserResponse GetUser(1: GetUserRequest req),
    CreateUserResponse CreateUser(1: CreateUserRequest req),
    ListUsersResponse ListUsers(1: ListUsersRequest req)
}
```

生成代码 / Generate code:
```bash
kitex -module your-module-name idl/user.thrift
```

## 3. 服务端实现 / Server Implementation

```go
// handler.go
package main

import (
    "context"
    "sync"
    "sync/atomic"
    
    "your-module/kitex_gen/api/user"
)

// UserServiceImpl 实现 UserService 接口
type UserServiceImpl struct {
    users  sync.Map
    nextID int64
}

func NewUserServiceImpl() *UserServiceImpl {
    return &UserServiceImpl{}
}

// GetUser 获取用户
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
    if u, ok := s.users.Load(req.Id); ok {
        return &user.GetUserResponse{
            User: u.(*user.User),
        }, nil
    }
    return &user.GetUserResponse{}, nil
}

// CreateUser 创建用户
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
    id := atomic.AddInt64(&s.nextID, 1)
    
    u := &user.User{
        Id:    id,
        Name:  req.Name,
        Email: req.Email,
        Age:   req.Age,
    }
    
    s.users.Store(id, u)
    
    return &user.CreateUserResponse{
        Id:      id,
        Success: true,
    }, nil
}

// ListUsers 列出用户
func (s *UserServiceImpl) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
    var users []*user.User
    s.users.Range(func(key, value interface{}) bool {
        users = append(users, value.(*user.User))
        return true
    })
    
    // 简单分页
    start := int(req.Page-1) * int(req.PageSize)
    end := start + int(req.PageSize)
    
    if start > len(users) {
        start = len(users)
    }
    if end > len(users) {
        end = len(users)
    }
    
    return &user.ListUsersResponse{
        Users: users[start:end],
        Total: int32(len(users)),
    }, nil
}
```

```go
// main.go (server)
package main

import (
    "log"
    "net"
    
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/server"
    "your-module/kitex_gen/api/user/userservice"
)

func main() {
    addr, _ := net.ResolveTCPAddr("tcp", ":8888")
    
    svr := userservice.NewServer(
        NewUserServiceImpl(),
        server.WithServiceAddr(addr),
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
            ServiceName: "user-service",
        }),
    )
    
    log.Println("Starting Kitex server on :8888")
    if err := svr.Run(); err != nil {
        log.Fatal("Server error:", err)
    }
}
```

## 4. 客户端实现 / Client Implementation

```go
// client/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/cloudwego/kitex/client"
    "github.com/cloudwego/kitex/client/callopt"
    "your-module/kitex_gen/api/user"
    "your-module/kitex_gen/api/user/userservice"
)

func main() {
    // 创建客户端 / Create client
    cli, err := userservice.NewClient(
        "user-service",
        client.WithHostPorts("127.0.0.1:8888"),
    )
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }
    
    ctx := context.Background()
    
    // 创建用户 / Create user
    createResp, err := cli.CreateUser(ctx, &user.CreateUserRequest{
        Name:  "Alice",
        Email: "alice@example.com",
        Age:   25,
    })
    if err != nil {
        log.Fatal("CreateUser error:", err)
    }
    fmt.Printf("Created user ID: %d\n", createResp.Id)
    
    // 获取用户 / Get user
    getResp, err := cli.GetUser(ctx, &user.GetUserRequest{
        Id: createResp.Id,
    })
    if err != nil {
        log.Fatal("GetUser error:", err)
    }
    fmt.Printf("Got user: %+v\n", getResp.User)
    
    // 带超时的调用 / Call with timeout
    ctx, cancel := context.WithTimeout(ctx, time.Second*3)
    defer cancel()
    
    listResp, err := cli.ListUsers(ctx, &user.ListUsersRequest{
        Page:     1,
        PageSize: 10,
    }, callopt.WithRPCTimeout(time.Second*2))
    if err != nil {
        log.Fatal("ListUsers error:", err)
    }
    fmt.Printf("Total users: %d\n", listResp.Total)
}
```

## 5. 中间件 / Middleware

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/cloudwego/kitex/pkg/endpoint"
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/server"
)

// 日志中间件 / Logging middleware
func LoggingMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
    return func(ctx context.Context, req, resp interface{}) error {
        start := time.Now()
        
        // 获取 RPC 信息 / Get RPC info
        ri := rpcinfo.GetRPCInfo(ctx)
        method := ri.Invocation().MethodName()
        
        log.Printf("[START] Method: %s", method)
        
        err := next(ctx, req, resp)
        
        duration := time.Since(start)
        if err != nil {
            log.Printf("[END] Method: %s, Duration: %v, Error: %v", method, duration, err)
        } else {
            log.Printf("[END] Method: %s, Duration: %v", method, duration)
        }
        
        return err
    }
}

// 恢复中间件 / Recovery middleware
func RecoveryMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
    return func(ctx context.Context, req, resp interface{}) (err error) {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Recovered from panic: %v", r)
                err = fmt.Errorf("internal error: %v", r)
            }
        }()
        return next(ctx, req, resp)
    }
}

// 使用中间件 / Using middleware
func main() {
    svr := userservice.NewServer(
        NewUserServiceImpl(),
        server.WithMiddleware(LoggingMiddleware),
        server.WithMiddleware(RecoveryMiddleware),
    )
    svr.Run()
}
```

## 6. 服务发现与负载均衡 / Service Discovery & Load Balancing

```go
package main

import (
    "github.com/cloudwego/kitex/client"
    "github.com/cloudwego/kitex/pkg/loadbalance"
    "github.com/cloudwego/kitex/pkg/retry"
    
    // Nacos 服务发现
    "github.com/kitex-contrib/registry-nacos/resolver"
    
    // Etcd 服务发现
    etcd "github.com/kitex-contrib/registry-etcd"
)

// 使用 Nacos / Using Nacos
func createClientWithNacos() {
    r, _ := resolver.NewDefaultNacosResolver()
    
    cli, _ := userservice.NewClient(
        "user-service",
        client.WithResolver(r),
        client.WithLoadBalancer(loadbalance.NewWeightedRandomBalancer()),
    )
    _ = cli
}

// 使用 Etcd / Using Etcd
func createClientWithEtcd() {
    r, _ := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
    
    cli, _ := userservice.NewClient(
        "user-service",
        client.WithResolver(r),
    )
    _ = cli
}

// 配置重试 / Configure retry
func createClientWithRetry() {
    // 配置重试策略 / Configure retry policy
    fp := retry.NewFailurePolicy()
    fp.WithMaxRetryTimes(3)
    fp.WithMaxDurationMS(5000)
    
    cli, _ := userservice.NewClient(
        "user-service",
        client.WithHostPorts("127.0.0.1:8888"),
        client.WithFailureRetry(fp),
    )
    _ = cli
}

// 配置熔断 / Configure circuit breaker
func createClientWithCircuitBreaker() {
    // Kitex 内置了熔断器
    // 可以通过 client.WithCircuitBreaker 配置
    
    cli, _ := userservice.NewClient(
        "user-service",
        client.WithHostPorts("127.0.0.1:8888"),
        // 自定义熔断配置
    )
    _ = cli
}
```

## 7. 链路追踪 / Tracing

```go
package main

import (
    "github.com/cloudwego/kitex/server"
    
    // OpenTelemetry
    "github.com/kitex-contrib/obs-opentelemetry/tracing"
)

func main() {
    // 初始化 tracer / Initialize tracer
    p := tracing.NewServerSuite()
    
    svr := userservice.NewServer(
        NewUserServiceImpl(),
        server.WithSuite(p),
    )
    
    svr.Run()
}

// 客户端配置 / Client configuration
func createClientWithTracing() {
    p := tracing.NewClientSuite()
    
    cli, _ := userservice.NewClient(
        "user-service",
        client.WithSuite(p),
    )
    _ = cli
}
```

## 8. 配置管理 / Configuration Management

```go
package main

import (
    "time"
    
    "github.com/cloudwego/kitex/client"
    "github.com/cloudwego/kitex/pkg/connpool"
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/server"
)

// 服务端配置 / Server configuration
func serverWithConfig() {
    svr := userservice.NewServer(
        NewUserServiceImpl(),
        server.WithServiceAddr(&net.TCPAddr{Port: 8888}),
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
            ServiceName: "user-service",
        }),
        server.WithReadWriteTimeout(time.Second * 5),
        server.WithExitWaitTime(time.Second * 3),
        server.WithMaxConnIdleTime(time.Minute * 10),
        server.WithMuxTransport(),  // 多路复用
    )
    _ = svr
}

// 客户端配置 / Client configuration
func clientWithConfig() {
    cli, _ := userservice.NewClient(
        "user-service",
        client.WithHostPorts("127.0.0.1:8888"),
        client.WithConnectTimeout(time.Second * 3),
        client.WithRPCTimeout(time.Second * 5),
        client.WithMuxConnection(2),  // 多路复用连接数
        client.WithLongConnection(connpool.IdleConfig{
            MaxIdlePerAddress: 10,
            MaxIdleGlobal:     100,
            MaxIdleTimeout:    time.Minute * 10,
        }),
    )
    _ = cli
}
```

## 9. gRPC 协议支持 / gRPC Protocol Support

```protobuf
// idl/user.proto
syntax = "proto3";

package api.user;

option go_package = "api/user";

message User {
    int64 id = 1;
    string name = 2;
    string email = 3;
    int32 age = 4;
}

message GetUserRequest {
    int64 id = 1;
}

message GetUserResponse {
    User user = 1;
}

service UserService {
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
}
```

生成代码 / Generate code:
```bash
kitex -module your-module-name -type protobuf idl/user.proto
```

## 10. 性能优化 / Performance Optimization

```go
package main

import (
    "github.com/cloudwego/kitex/pkg/transmeta"
    "github.com/cloudwego/kitex/server"
    "github.com/cloudwego/kitex/transport"
)

func optimizedServer() {
    svr := userservice.NewServer(
        NewUserServiceImpl(),
        // 使用 netpoll (Linux)
        // 默认已启用
        
        // 多路复用传输
        server.WithMuxTransport(),
        
        // 调整连接池
        server.WithLimit(&limit.Option{
            MaxConnections: 10000,
            MaxQPS:         100000,
        }),
        
        // 元信息传输
        server.WithMetaHandler(transmeta.ServerTTHeaderHandler),
    )
    _ = svr
}

// 最佳实践 / Best practices:
// 1. 使用多路复用传输减少连接数
// 2. 合理配置连接池大小
// 3. 启用 gzip 压缩 (对于大消息)
// 4. 使用异步调用提高吞吐量
// 5. 合理设置超时时间
```
