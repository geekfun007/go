# Apache Thrift - RPC 框架 / RPC Framework

## 1. Thrift 概述 / Thrift Overview

Apache Thrift 是一个可扩展的跨语言服务开发框架，用于构建高效、可扩展的服务。

```
安装 / Installation:
# 安装 Thrift 编译器 / Install Thrift compiler
# macOS
brew install thrift

# Ubuntu
apt-get install thrift-compiler

# Go Thrift 库 / Go Thrift library
go get github.com/apache/thrift/lib/go/thrift
```

## 2. 定义 Thrift IDL / Define Thrift IDL

```thrift
// user.thrift

namespace go user

// 枚举类型 / Enum type
enum UserStatus {
    ACTIVE = 1,
    INACTIVE = 2,
    BANNED = 3
}

// 结构体 / Struct
struct User {
    1: required i64 id,
    2: required string name,
    3: optional string email,
    4: i32 age,
    5: UserStatus status = UserStatus.ACTIVE,
    6: list<string> tags,
    7: map<string, string> extra
}

// 请求结构 / Request struct
struct GetUserRequest {
    1: required i64 user_id
}

struct CreateUserRequest {
    1: required string name,
    2: optional string email,
    3: i32 age
}

struct UpdateUserRequest {
    1: required i64 user_id,
    2: optional string name,
    3: optional string email,
    4: optional i32 age
}

// 响应结构 / Response struct
struct GetUserResponse {
    1: User user,
    2: string message
}

struct CreateUserResponse {
    1: i64 user_id,
    2: bool success,
    3: string message
}

struct ListUsersResponse {
    1: list<User> users,
    2: i32 total
}

// 异常 / Exception
exception UserNotFoundException {
    1: i64 user_id,
    2: string message
}

exception ValidationException {
    1: string field,
    2: string message
}

// 服务定义 / Service definition
service UserService {
    // 获取用户 / Get user
    GetUserResponse getUser(1: GetUserRequest req) throws (
        1: UserNotFoundException notFound
    ),
    
    // 创建用户 / Create user
    CreateUserResponse createUser(1: CreateUserRequest req) throws (
        1: ValidationException validation
    ),
    
    // 更新用户 / Update user
    void updateUser(1: UpdateUserRequest req) throws (
        1: UserNotFoundException notFound,
        2: ValidationException validation
    ),
    
    // 删除用户 / Delete user
    void deleteUser(1: i64 userId) throws (
        1: UserNotFoundException notFound
    ),
    
    // 列出用户 / List users
    ListUsersResponse listUsers(1: i32 offset, 2: i32 limit),
    
    // 单向方法 (无响应) / Oneway method (no response)
    oneway void ping()
}
```

生成代码 / Generate code:
```bash
thrift -r --gen go user.thrift
```

## 3. 服务端实现 / Server Implementation

```go
package main

import (
    "context"
    "fmt"
    "log"
    "sync"
    
    "github.com/apache/thrift/lib/go/thrift"
    "your-module/gen-go/user"  // 生成的代码
)

// 实现 UserService 接口 / Implement UserService interface
type UserServiceHandler struct {
    users  map[int64]*user.User
    nextID int64
    mu     sync.RWMutex
}

func NewUserServiceHandler() *UserServiceHandler {
    return &UserServiceHandler{
        users:  make(map[int64]*user.User),
        nextID: 1,
    }
}

// GetUser 获取用户 / Get user
func (h *UserServiceHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    u, ok := h.users[req.UserID]
    if !ok {
        return nil, &user.UserNotFoundException{
            UserID:  req.UserID,
            Message: fmt.Sprintf("User %d not found", req.UserID),
        }
    }
    
    return &user.GetUserResponse{
        User:    u,
        Message: "success",
    }, nil
}

// CreateUser 创建用户 / Create user
func (h *UserServiceHandler) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
    // 验证 / Validation
    if req.Name == "" {
        return nil, &user.ValidationException{
            Field:   "name",
            Message: "Name is required",
        }
    }
    
    h.mu.Lock()
    defer h.mu.Unlock()
    
    id := h.nextID
    h.nextID++
    
    h.users[id] = &user.User{
        ID:     id,
        Name:   req.Name,
        Email:  req.Email,
        Age:    req.Age,
        Status: user.UserStatus_ACTIVE,
    }
    
    return &user.CreateUserResponse{
        UserID:  id,
        Success: true,
        Message: "User created successfully",
    }, nil
}

// UpdateUser 更新用户 / Update user
func (h *UserServiceHandler) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) error {
    h.mu.Lock()
    defer h.mu.Unlock()
    
    u, ok := h.users[req.UserID]
    if !ok {
        return &user.UserNotFoundException{
            UserID:  req.UserID,
            Message: fmt.Sprintf("User %d not found", req.UserID),
        }
    }
    
    if req.Name != nil {
        u.Name = *req.Name
    }
    if req.Email != nil {
        u.Email = req.Email
    }
    if req.Age != nil {
        u.Age = *req.Age
    }
    
    return nil
}

// DeleteUser 删除用户 / Delete user
func (h *UserServiceHandler) DeleteUser(ctx context.Context, userID int64) error {
    h.mu.Lock()
    defer h.mu.Unlock()
    
    if _, ok := h.users[userID]; !ok {
        return &user.UserNotFoundException{
            UserID:  userID,
            Message: fmt.Sprintf("User %d not found", userID),
        }
    }
    
    delete(h.users, userID)
    return nil
}

// ListUsers 列出用户 / List users
func (h *UserServiceHandler) ListUsers(ctx context.Context, offset int32, limit int32) (*user.ListUsersResponse, error) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    users := make([]*user.User, 0)
    for _, u := range h.users {
        users = append(users, u)
    }
    
    // 分页 / Pagination
    start := int(offset)
    end := start + int(limit)
    if start > len(users) {
        start = len(users)
    }
    if end > len(users) {
        end = len(users)
    }
    
    return &user.ListUsersResponse{
        Users: users[start:end],
        Total: int32(len(h.users)),
    }, nil
}

// Ping 心跳 / Heartbeat
func (h *UserServiceHandler) Ping(ctx context.Context) error {
    log.Println("Ping received")
    return nil
}

// 启动服务器 / Start server
func main() {
    handler := NewUserServiceHandler()
    processor := user.NewUserServiceProcessor(handler)
    
    // 传输层配置 / Transport configuration
    serverTransport, err := thrift.NewTServerSocket(":9090")
    if err != nil {
        log.Fatal("Error creating server socket:", err)
    }
    
    // 协议工厂 / Protocol factory
    protocolFactory := thrift.NewTBinaryProtocolFactoryConf(nil)
    
    // 传输工厂 / Transport factory
    transportFactory := thrift.NewTTransportFactory()
    
    // 可选: 使用带缓冲的传输 / Optional: Use buffered transport
    // transportFactory := thrift.NewTBufferedTransportFactory(8192)
    
    // 可选: 使用帧传输 / Optional: Use framed transport
    // transportFactory := thrift.NewTFramedTransportFactoryConf(nil)
    
    // 创建服务器 / Create server
    server := thrift.NewTSimpleServer4(
        processor,
        serverTransport,
        transportFactory,
        protocolFactory,
    )
    
    log.Println("Starting Thrift server on :9090")
    if err := server.Serve(); err != nil {
        log.Fatal("Error starting server:", err)
    }
}
```

## 4. 客户端实现 / Client Implementation

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/apache/thrift/lib/go/thrift"
    "your-module/gen-go/user"
)

func main() {
    // 配置 / Configuration
    cfg := &thrift.TConfiguration{}
    
    // 创建传输 / Create transport
    transport := thrift.NewTSocketConf("localhost:9090", cfg)
    
    // 使用缓冲传输 / Use buffered transport
    bufferedTransport := thrift.NewTBufferedTransport(transport, 8192)
    
    // 打开连接 / Open connection
    if err := bufferedTransport.Open(); err != nil {
        log.Fatal("Error opening transport:", err)
    }
    defer bufferedTransport.Close()
    
    // 创建协议 / Create protocol
    protocolFactory := thrift.NewTBinaryProtocolFactoryConf(cfg)
    protocol := protocolFactory.GetProtocol(bufferedTransport)
    
    // 创建客户端 / Create client
    client := user.NewUserServiceClient(thrift.NewTStandardClient(protocol, protocol))
    
    ctx := context.Background()
    
    // 创建用户 / Create user
    createResp, err := client.CreateUser(ctx, &user.CreateUserRequest{
        Name:  "Alice",
        Email: thrift.StringPtr("alice@example.com"),
        Age:   25,
    })
    if err != nil {
        log.Fatal("Error creating user:", err)
    }
    fmt.Printf("Created user with ID: %d\n", createResp.UserID)
    
    // 获取用户 / Get user
    getResp, err := client.GetUser(ctx, &user.GetUserRequest{
        UserID: createResp.UserID,
    })
    if err != nil {
        switch e := err.(type) {
        case *user.UserNotFoundException:
            log.Printf("User not found: %s\n", e.Message)
        default:
            log.Fatal("Error getting user:", err)
        }
        return
    }
    fmt.Printf("Got user: %+v\n", getResp.User)
    
    // 更新用户 / Update user
    newName := "Alice Updated"
    err = client.UpdateUser(ctx, &user.UpdateUserRequest{
        UserID: createResp.UserID,
        Name:   &newName,
    })
    if err != nil {
        log.Fatal("Error updating user:", err)
    }
    fmt.Println("User updated successfully")
    
    // 列出用户 / List users
    listResp, err := client.ListUsers(ctx, 0, 10)
    if err != nil {
        log.Fatal("Error listing users:", err)
    }
    fmt.Printf("Total users: %d\n", listResp.Total)
    for _, u := range listResp.Users {
        fmt.Printf("  - %d: %s\n", u.ID, u.Name)
    }
    
    // 发送心跳 / Send ping
    client.Ping(ctx)
}

// 连接池实现 / Connection pool implementation
type ThriftClientPool struct {
    clients chan *user.UserServiceClient
    factory func() (*user.UserServiceClient, error)
    maxSize int
}

func NewThriftClientPool(maxSize int, addr string) (*ThriftClientPool, error) {
    pool := &ThriftClientPool{
        clients: make(chan *user.UserServiceClient, maxSize),
        maxSize: maxSize,
        factory: func() (*user.UserServiceClient, error) {
            cfg := &thrift.TConfiguration{}
            transport := thrift.NewTSocketConf(addr, cfg)
            bufferedTransport := thrift.NewTBufferedTransport(transport, 8192)
            
            if err := bufferedTransport.Open(); err != nil {
                return nil, err
            }
            
            protocolFactory := thrift.NewTBinaryProtocolFactoryConf(cfg)
            protocol := protocolFactory.GetProtocol(bufferedTransport)
            
            return user.NewUserServiceClient(
                thrift.NewTStandardClient(protocol, protocol),
            ), nil
        },
    }
    
    // 预创建连接 / Pre-create connections
    for i := 0; i < maxSize; i++ {
        client, err := pool.factory()
        if err != nil {
            return nil, err
        }
        pool.clients <- client
    }
    
    return pool, nil
}

func (p *ThriftClientPool) Get() *user.UserServiceClient {
    return <-p.clients
}

func (p *ThriftClientPool) Put(client *user.UserServiceClient) {
    p.clients <- client
}
```

## 5. 高级配置 / Advanced Configuration

```go
package main

import (
    "crypto/tls"
    "time"
    
    "github.com/apache/thrift/lib/go/thrift"
)

// TLS 配置 / TLS configuration
func createTLSTransport(addr string) (*thrift.TSocket, error) {
    tlsConfig := &tls.Config{
        InsecureSkipVerify: true,  // 生产环境应该设置正确的证书
    }
    
    cfg := &thrift.TConfiguration{
        ConnectTimeout: time.Second * 10,
        SocketTimeout:  time.Second * 10,
    }
    
    transport := thrift.NewTSSLSocketConf(addr, tlsConfig, cfg)
    return transport, nil
}

// 不同协议 / Different protocols
func protocols() {
    cfg := &thrift.TConfiguration{}
    
    // 二进制协议 (默认) / Binary protocol (default)
    _ = thrift.NewTBinaryProtocolFactoryConf(cfg)
    
    // 紧凑协议 (更小的消息体) / Compact protocol (smaller messages)
    _ = thrift.NewTCompactProtocolFactoryConf(cfg)
    
    // JSON 协议 (可读性好) / JSON protocol (readable)
    _ = thrift.NewTJSONProtocolFactory()
}

// 不同传输方式 / Different transports
func transports() {
    // 缓冲传输 / Buffered transport
    _ = thrift.NewTBufferedTransportFactory(8192)
    
    // 帧传输 / Framed transport
    _ = thrift.NewTFramedTransportFactoryConf(nil)
    
    // HTTP 传输 / HTTP transport
    // 需要配合 HTTP server 使用
}
```
