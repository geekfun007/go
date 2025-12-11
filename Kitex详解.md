# Kitex 详解与实战

## 目录
- [简介](#简介)
- [Kitex CLI 命令行详解](#kitex-cli-命令行详解)
- [Context 详解](#context-详解)
- [服务定义与代码生成](#服务定义与代码生成)
- [服务端开发](#服务端开发)
- [客户端开发](#客户端开发)
- [中间件](#中间件)
- [服务治理](#服务治理)
- [性能优化](#性能优化)
- [实战案例](#实战案例)

---

## 简介

Kitex 是 CloudWeGo 开源的高性能、强可扩展的 Go RPC 框架。它支持 Thrift、Protobuf、gRPC 等多种协议，内置丰富的服务治理能力。

### 主要特性

- **高性能**：基于自研的高性能网络库 Netpoll，性能优于 gRPC
- **扩展性强**：提供大量的扩展接口，支持自定义协议、传输、编解码等
- **多协议**：支持 Thrift、Protobuf、gRPC
- **多消息类型**：支持 PingPong、Oneway、Streaming（双向流、客户端流、服务端流）
- **服务治理**：集成服务注册发现、负载均衡、熔断、限流等
- **代码生成**：提供强大的代码生成工具
- **可观测性**：完整的监控、日志、链路追踪支持

### 安装

```bash
# 安装 Kitex
go get github.com/cloudwego/kitex

# 安装 Kitex 命令行工具
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest

# 安装 Thriftgo（如果使用 Thrift IDL）
go install github.com/cloudwego/thriftgo@latest

# 验证安装
kitex --version
```

---

## Kitex CLI 命令行详解

Kitex 提供了 `kitex` 命令行工具，用于生成 RPC 代码。

### kitex 命令总览

```bash
kitex --help
```

输出：
```
NAME:
   kitex - A code generation tool for Kitex framework

USAGE:
   kitex [global options] command [command options] IDL

VERSION:
   v0.x.x

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --module value, -m value         specify the Go module name (required)
   --service value, -s value        specify the service name
   --type value                     specify the generated code type: thrift or protobuf (default: "thrift")
   --includes value, -I value       add IDL search path (can be specified multiple times)
   --thrift value                   thrift compiler path (default: use $PATH)
   --protoc value                   protoc compiler path (default: use $PATH)
   --proto-search-path value        protobuf IDL search path (can be specified multiple times)
   --combine-service                combine services into one client
   --copy-idl                       copy IDL to output directory
   --gen-path value                 specify the code generation path (default: current directory)
   --verbose                        verbose mode
   --help, -h                       show help
   --version, -v                    print the version
```

---

### kitex 基本用法

#### 1. 从 Thrift IDL 生成代码

```bash
# 基本用法
kitex -module github.com/myuser/myproject idl/api.thrift

# 指定服务名
kitex -module github.com/myuser/myproject -service myservice idl/api.thrift

# 指定生成路径
kitex -module github.com/myuser/myproject -gen-path ./gen idl/api.thrift

# 指定 include 路径
kitex -module github.com/myuser/myproject -I idl -I ../common idl/api.thrift

# 复制 IDL 文件到输出目录
kitex -module github.com/myuser/myproject -copy-idl idl/api.thrift
```

#### 2. 从 Protobuf IDL 生成代码

```bash
# 基本用法
kitex -module github.com/myuser/myproject -type protobuf idl/api.proto

# 指定 proto 搜索路径
kitex -module github.com/myuser/myproject \
  -type protobuf \
  -proto-search-path idl \
  idl/api.proto

# 组合服务（多个 service 生成一个 client）
kitex -module github.com/myuser/myproject \
  -type protobuf \
  -combine-service \
  idl/api.proto
```

---

### 完整参数详解

#### 必需参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `--module, -m` | Go module 名称（必需） | `-module github.com/myuser/myproject` |
| `IDL` | IDL 文件路径（位置参数） | `idl/api.thrift` |

#### 服务相关

| 参数 | 说明 | 示例 |
|------|------|------|
| `--service, -s` | 指定服务名称 | `-service userservice` |
| `--type` | IDL 类型：thrift 或 protobuf | `-type protobuf` |
| `--combine-service` | 合并多个服务到一个 client | `--combine-service` |

#### 路径相关

| 参数 | 说明 | 示例 |
|------|------|------|
| `--gen-path` | 代码生成路径 | `-gen-path ./gen` |
| `--includes, -I` | IDL 搜索路径（可多次指定） | `-I idl -I ../common` |
| `--proto-search-path` | Protobuf IDL 搜索路径 | `-proto-search-path idl` |

#### 工具路径

| 参数 | 说明 | 示例 |
|------|------|------|
| `--thrift` | thriftgo 编译器路径 | `--thrift /path/to/thriftgo` |
| `--protoc` | protoc 编译器路径 | `--protoc /path/to/protoc` |

#### 其他选项

| 参数 | 说明 | 示例 |
|------|------|------|
| `--copy-idl` | 复制 IDL 到输出目录 | `--copy-idl` |
| `--verbose` | 详细输出模式 | `--verbose` |
| `--help, -h` | 显示帮助信息 | `--help` |
| `--version, -v` | 显示版本信息 | `--version` |

---

### 实战示例

#### 示例 1：创建 Thrift 微服务

**步骤 1：创建 IDL 文件**

**idl/user.thrift**
```thrift
namespace go user

// 用户信息
struct User {
    1: required i64 id
    2: required string username
    3: required string email
    4: optional i32 age
    5: optional string created_at
}

// 创建用户请求
struct CreateUserRequest {
    1: required string username
    2: required string email
    3: optional i32 age
}

// 创建用户响应
struct CreateUserResponse {
    1: required i32 code
    2: required string message
    3: optional User user
}

// 获取用户请求
struct GetUserRequest {
    1: required i64 user_id
}

// 获取用户响应
struct GetUserResponse {
    1: required i32 code
    2: required string message
    3: optional User user
}

// 用户服务
service UserService {
    CreateUserResponse CreateUser(1: CreateUserRequest req)
    GetUserResponse GetUser(1: GetUserRequest req)
}
```

**步骤 2：生成代码**

```bash
# 创建项目目录
mkdir user-service
cd user-service

# 初始化 Go module
go mod init github.com/myuser/user-service

# 生成代码
kitex -module github.com/myuser/user-service \
    -service userservice \
    idl/user.thrift

# 下载依赖
go mod tidy
```

**生成的目录结构：**
```
user-service/
├── kitex_gen/              # 生成的基础代码
│   └── user/
│       ├── user.go         # 数据结构
│       ├── userservice/    # 服务代码
│       │   ├── client.go   # 客户端
│       │   ├── server.go   # 服务端
│       │   ├── userservice.go
│       │   └── invoker.go
│       └── k-consts.go
├── handler.go              # 业务逻辑处理器（需要实现）
├── main.go                 # 服务入口
├── build.sh                # 构建脚本
├── script/                 # 辅助脚本
│   └── bootstrap.sh
├── idl/
│   └── user.thrift
├── go.mod
└── go.sum
```

#### 示例 2：创建 Protobuf 微服务

**步骤 1：创建 IDL 文件**

**idl/user.proto**
```protobuf
syntax = "proto3";

package user;

option go_package = "github.com/myuser/user-service/kitex_gen/user";

// 用户信息
message User {
  int64 id = 1;
  string username = 2;
  string email = 3;
  int32 age = 4;
  string created_at = 5;
}

// 创建用户请求
message CreateUserRequest {
  string username = 1;
  string email = 2;
  int32 age = 3;
}

// 创建用户响应
message CreateUserResponse {
  int32 code = 1;
  string message = 2;
  User user = 3;
}

// 获取用户请求
message GetUserRequest {
  int64 user_id = 1;
}

// 获取用户响应
message GetUserResponse {
  int32 code = 1;
  string message = 2;
  User user = 3;
}

// 用户服务
service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}
```

**步骤 2：生成代码**

```bash
# 生成代码
kitex -module github.com/myuser/user-service \
    -service userservice \
    -type protobuf \
    idl/user.proto

# 下载依赖
go mod tidy
```

#### 示例 3：多服务项目

**项目结构：**
```
myproject/
├── idl/
│   ├── common.thrift       # 公共定义
│   ├── user.thrift         # 用户服务
│   └── post.thrift         # 文章服务
├── services/
│   ├── user/               # 用户服务
│   └── post/               # 文章服务
└── api/                    # API 网关
```

**common.thrift**
```thrift
namespace go common

struct BaseResponse {
    1: required i32 code
    2: required string message
}
```

**user.thrift**
```thrift
include "common.thrift"

namespace go user

struct User {
    1: required i64 id
    2: required string username
}

struct GetUserRequest {
    1: required i64 user_id
}

struct GetUserResponse {
    1: required common.BaseResponse base
    2: optional User user
}

service UserService {
    GetUserResponse GetUser(1: GetUserRequest req)
}
```

**生成用户服务代码：**
```bash
cd services/user

kitex -module github.com/myuser/myproject/services/user \
    -service userservice \
    -I ../../idl \
    ../../idl/user.thrift

go mod tidy
```

#### 示例 4：更新已有服务

当 IDL 文件变更后，重新运行 kitex 命令更新代码：

```bash
# 更新代码（会保留 handler.go 中的业务逻辑）
kitex -module github.com/myuser/user-service \
    -service userservice \
    idl/user.thrift
```

⚠️ **注意**：
- `handler.go` 中已实现的业务逻辑会被保留
- `main.go`、`kitex_gen/` 等生成的代码会被更新
- 建议使用版本控制系统管理代码

---

### 常用命令模板

#### 快速生成服务脚本

**generate_service.sh**
```bash
#!/bin/bash

# 配置
SERVICE_NAME=$1
MODULE_NAME="github.com/mycompany/${SERVICE_NAME}"
IDL_FILE="idl/${SERVICE_NAME}.thrift"

if [ -z "$SERVICE_NAME" ]; then
    echo "Usage: $0 <service-name>"
    exit 1
fi

# 创建目录
mkdir -p ${SERVICE_NAME}/idl

# 生成代码
cd ${SERVICE_NAME}

kitex -module ${MODULE_NAME} \
    -service ${SERVICE_NAME} \
    -copy-idl \
    ../${IDL_FILE}

# 下载依赖
go mod tidy

echo "Service ${SERVICE_NAME} generated successfully!"
echo "Next steps:"
echo "1. Implement business logic in handler.go"
echo "2. Run: go run ."
```

**使用方法：**
```bash
chmod +x generate_service.sh
./generate_service.sh user-service
```

---

## Context 详解

Kitex 使用标准的 `context.Context` 进行请求上下文传递，同时提供了 `metainfo` 包用于元信息传递。

### context.Context 使用

```go
package main

import (
    "context"
    "time"
)

// 1. 超时控制
func WithTimeout() {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    // 调用 RPC 方法
    resp, err := client.GetUser(ctx, req)
    if err != nil {
        // 处理超时错误
        if ctx.Err() == context.DeadlineExceeded {
            log.Println("request timeout")
        }
    }
}

// 2. 取消控制
func WithCancel() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go func() {
        // 某个条件下取消请求
        time.Sleep(1 * time.Second)
        cancel()
    }()
    
    resp, err := client.GetUser(ctx, req)
    if err != nil {
        if ctx.Err() == context.Canceled {
            log.Println("request canceled")
        }
    }
}

// 3. 值传递
func WithValue() {
    ctx := context.WithValue(context.Background(), "request_id", "abc123")
    
    resp, err := client.GetUser(ctx, req)
}
```

### metainfo - 元信息传递

Kitex 提供 `metainfo` 包用于在客户端和服务端之间传递元信息。

#### 客户端设置元信息

```go
package main

import (
    "context"
    
    "github.com/cloudwego/kitex/pkg/metainfo"
)

func ClientSetMetaInfo() {
    ctx := context.Background()
    
    // 方式 1：使用 WithValue 设置单个值
    ctx = metainfo.WithValue(ctx, "user_id", "123")
    ctx = metainfo.WithValue(ctx, "request_id", "abc-def-ghi")
    
    // 方式 2：使用 WithPersistentValue 设置持久化值（会传递到下游服务）
    ctx = metainfo.WithPersistentValue(ctx, "trace_id", "trace-123")
    
    // 方式 3：批量设置
    ctx = metainfo.WithValues(ctx, map[string]string{
        "user_id":    "123",
        "request_id": "abc-def",
        "client_ip":  "192.168.1.1",
    })
    
    // 调用 RPC
    resp, err := client.GetUser(ctx, req)
}
```

#### 服务端获取元信息

```go
package main

import (
    "context"
    
    "github.com/cloudwego/kitex/pkg/metainfo"
)

// UserServiceImpl 实现 UserService 接口
type UserServiceImpl struct{}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
    // 方式 1：获取单个值
    userID, ok := metainfo.GetValue(ctx, "user_id")
    if ok {
        log.Printf("User ID: %s", userID)
    }
    
    // 方式 2：获取持久化值
    traceID, ok := metainfo.GetPersistentValue(ctx, "trace_id")
    if ok {
        log.Printf("Trace ID: %s", traceID)
    }
    
    // 方式 3：获取所有值
    allValues := metainfo.GetAllValues(ctx)
    for k, v := range allValues {
        log.Printf("%s: %s", k, v)
    }
    
    // 业务逻辑
    return &user.GetUserResponse{
        Code:    0,
        Message: "success",
        User: &user.User{
            Id:       req.UserId,
            Username: "zhangsan",
            Email:    "zhangsan@example.com",
        },
    }, nil
}
```

### 链路传递

```go
// 服务 A 调用服务 B，服务 B 调用服务 C

// 服务 A（客户端）
func ServiceA() {
    ctx := context.Background()
    
    // 设置持久化元信息（会自动传递到下游）
    ctx = metainfo.WithPersistentValue(ctx, "trace_id", "trace-123")
    ctx = metainfo.WithPersistentValue(ctx, "user_id", "user-456")
    
    // 调用服务 B
    respB, err := clientB.CallServiceB(ctx, reqB)
}

// 服务 B（服务端 + 客户端）
func (s *ServiceBImpl) CallServiceB(ctx context.Context, req *RequestB) (*ResponseB, error) {
    // 获取上游传递的元信息
    traceID, _ := metainfo.GetPersistentValue(ctx, "trace_id")
    log.Printf("Service B received trace_id: %s", traceID)
    
    // 继续传递到服务 C（ctx 会自动携带持久化元信息）
    respC, err := clientC.CallServiceC(ctx, reqC)
    
    return &ResponseB{}, nil
}

// 服务 C（服务端）
func (s *ServiceCImpl) CallServiceC(ctx context.Context, req *RequestC) (*ResponseC, error) {
    // 仍然可以获取到服务 A 设置的元信息
    traceID, _ := metainfo.GetPersistentValue(ctx, "trace_id")
    userID, _ := metainfo.GetPersistentValue(ctx, "user_id")
    
    log.Printf("Service C received trace_id: %s, user_id: %s", traceID, userID)
    
    return &ResponseC{}, nil
}
```

### 获取请求信息

```go
import (
    "github.com/cloudwego/kitex/pkg/rpcinfo"
)

func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
    // 获取 RPC 信息
    ri := rpcinfo.GetRPCInfo(ctx)
    
    // 获取调用方信息
    from := ri.From()
    fromService := from.ServiceName()
    fromMethod := from.Method()
    fromAddress := from.Address()
    
    log.Printf("From service: %s, method: %s, address: %s", 
        fromService, fromMethod, fromAddress)
    
    // 获取被调用方信息
    to := ri.To()
    toService := to.ServiceName()
    toMethod := to.Method()
    
    log.Printf("To service: %s, method: %s", toService, toMethod)
    
    // 获取调用配置
    cfg := ri.Config()
    timeout := cfg.RPCTimeout()
    log.Printf("RPC timeout: %v", timeout)
    
    // 业务逻辑
    return &user.GetUserResponse{}, nil
}
```

---

## 服务定义与代码生成

### Thrift IDL 编写规范

#### 1. 基本类型

```thrift
namespace go example

// 基本类型
struct BasicTypes {
    1: bool bool_field           // 布尔
    2: byte byte_field           // 字节
    3: i16 i16_field             // 16位整数
    4: i32 i32_field             // 32位整数
    5: i64 i64_field             // 64位整数
    6: double double_field       // 双精度浮点数
    7: string string_field       // 字符串
    8: binary binary_field       // 二进制数据
}
```

#### 2. 容器类型

```thrift
// 容器类型
struct ContainerTypes {
    1: list<string> string_list          // 列表
    2: set<i32> int_set                  // 集合
    3: map<string, i32> string_int_map   // 映射
}
```

#### 3. 字段修饰符

```thrift
struct FieldModifiers {
    1: required string username    // 必填字段
    2: optional string email       // 可选字段
    3: string phone                // 默认（等同于 optional）
}
```

#### 4. 默认值

```thrift
struct DefaultValues {
    1: i32 status = 0
    2: string message = "success"
    3: bool active = true
}
```

#### 5. 枚举

```thrift
enum Gender {
    UNKNOWN = 0
    MALE = 1
    FEMALE = 2
}

struct User {
    1: i64 id
    2: string username
    3: Gender gender
}
```

#### 6. 常量

```thrift
const i32 MAX_USER_AGE = 150
const string DEFAULT_USERNAME = "guest"
const list<string> ALLOWED_ROLES = ["admin", "user", "guest"]
```

#### 7. 异常

```thrift
exception UserNotFoundError {
    1: i32 code
    2: string message
}

exception ValidationError {
    1: i32 code
    2: string message
    3: map<string, string> errors
}

service UserService {
    User GetUser(1: i64 user_id) throws (1: UserNotFoundError notFound)
}
```

#### 8. 服务继承

```thrift
// 基础服务
service BaseService {
    string Ping()
}

// 继承基础服务
service UserService extends BaseService {
    User GetUser(1: i64 user_id)
}
```

#### 9. Include 和 Namespace

```thrift
// common.thrift
namespace go common

struct BaseResponse {
    1: i32 code
    2: string message
}
```

```thrift
// user.thrift
include "common.thrift"

namespace go user

struct GetUserResponse {
    1: common.BaseResponse base
    2: optional User user
}
```

### Protobuf IDL 编写规范

#### 1. 基本语法

```protobuf
syntax = "proto3";

package user;

option go_package = "github.com/myuser/myproject/kitex_gen/user";

// 用户信息
message User {
  int64 id = 1;
  string username = 2;
  string email = 3;
  int32 age = 4;
}
```

#### 2. 字段类型

```protobuf
message FieldTypes {
  // 数值类型
  int32 int32_field = 1;
  int64 int64_field = 2;
  uint32 uint32_field = 3;
  uint64 uint64_field = 4;
  float float_field = 5;
  double double_field = 6;
  
  // 布尔和字符串
  bool bool_field = 7;
  string string_field = 8;
  bytes bytes_field = 9;
  
  // 重复字段（列表）
  repeated string tags = 10;
  
  // 映射
  map<string, int32> scores = 11;
}
```

#### 3. 枚举

```protobuf
enum Gender {
  GENDER_UNKNOWN = 0;
  GENDER_MALE = 1;
  GENDER_FEMALE = 2;
}

message User {
  int64 id = 1;
  string username = 2;
  Gender gender = 3;
}
```

#### 4. 嵌套消息

```protobuf
message User {
  int64 id = 1;
  string username = 2;
  
  message Address {
    string country = 1;
    string city = 2;
    string street = 3;
  }
  
  Address address = 3;
}
```

#### 5. Oneof

```protobuf
message SearchRequest {
  string query = 1;
  
  oneof filter {
    int32 user_id = 2;
    string username = 3;
    string email = 4;
  }
}
```

#### 6. 服务定义

```protobuf
service UserService {
  // 简单 RPC
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  
  // 服务端流
  rpc ListUsers(ListUsersRequest) returns (stream User);
  
  // 客户端流
  rpc CreateUsers(stream CreateUserRequest) returns (CreateUsersResponse);
  
  // 双向流
  rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}
```

#### 7. Import

```protobuf
// common.proto
syntax = "proto3";

package common;

message BaseResponse {
  int32 code = 1;
  string message = 2;
}
```

```protobuf
// user.proto
syntax = "proto3";

package user;

import "common.proto";

message GetUserResponse {
  common.BaseResponse base = 1;
  User user = 2;
}
```

---

## 服务端开发

### 实现服务接口

#### Thrift 服务实现

```go
package main

import (
    "context"
    "errors"
    "time"
    
    "github.com/myuser/user-service/kitex_gen/user"
)

// UserServiceImpl 实现 UserService 接口
type UserServiceImpl struct {
    // 可以注入依赖
    userRepo *UserRepository
    cache    *Cache
}

// NewUserService 创建服务实例
func NewUserService(repo *UserRepository, cache *Cache) *UserServiceImpl {
    return &UserServiceImpl{
        userRepo: repo,
        cache:    cache,
    }
}

// CreateUser 实现创建用户方法
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
    // 1. 参数验证
    if req.Username == "" {
        return &user.CreateUserResponse{
            Code:    400,
            Message: "username is required",
        }, nil
    }
    
    if req.Email == "" {
        return &user.CreateUserResponse{
            Code:    400,
            Message: "email is required",
        }, nil
    }
    
    // 2. 业务逻辑
    newUser := &User{
        Username:  req.Username,
        Email:     req.Email,
        Age:       req.Age,
        CreatedAt: time.Now(),
    }
    
    // 3. 数据库操作
    if err := s.userRepo.Create(newUser); err != nil {
        return &user.CreateUserResponse{
            Code:    500,
            Message: "failed to create user",
        }, err
    }
    
    // 4. 缓存操作
    s.cache.SetUser(newUser.ID, newUser)
    
    // 5. 返回响应
    return &user.CreateUserResponse{
        Code:    0,
        Message: "success",
        User: &user.User{
            Id:        newUser.ID,
            Username:  newUser.Username,
            Email:     newUser.Email,
            Age:       newUser.Age,
            CreatedAt: newUser.CreatedAt.Format(time.RFC3339),
        },
    }, nil
}

// GetUser 实现获取用户方法
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
    // 1. 参数验证
    if req.UserId <= 0 {
        return &user.GetUserResponse{
            Code:    400,
            Message: "invalid user id",
        }, nil
    }
    
    // 2. 先查缓存
    if cachedUser, ok := s.cache.GetUser(req.UserId); ok {
        return &user.GetUserResponse{
            Code:    0,
            Message: "success",
            User:    convertToThriftUser(cachedUser),
        }, nil
    }
    
    // 3. 查数据库
    dbUser, err := s.userRepo.GetByID(req.UserId)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return &user.GetUserResponse{
                Code:    404,
                Message: "user not found",
            }, nil
        }
        return &user.GetUserResponse{
            Code:    500,
            Message: "internal error",
        }, err
    }
    
    // 4. 更新缓存
    s.cache.SetUser(dbUser.ID, dbUser)
    
    // 5. 返回响应
    return &user.GetUserResponse{
        Code:    0,
        Message: "success",
        User:    convertToThriftUser(dbUser),
    }, nil
}

func convertToThriftUser(u *User) *user.User {
    return &user.User{
        Id:        u.ID,
        Username:  u.Username,
        Email:     u.Email,
        Age:       u.Age,
        CreatedAt: u.CreatedAt.Format(time.RFC3339),
    }
}
```

### 启动服务器

#### 基本启动

```go
package main

import (
    "log"
    
    "github.com/cloudwego/kitex/server"
    "github.com/myuser/user-service/kitex_gen/user/userservice"
)

func main() {
    // 创建服务实现
    impl := NewUserService(userRepo, cache)
    
    // 创建服务器
    svr := userservice.NewServer(impl)
    
    // 启动服务器
    err := svr.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

#### 高级配置

```go
package main

import (
    "log"
    "net"
    "time"
    
    "github.com/cloudwego/kitex/pkg/klog"
    "github.com/cloudwego/kitex/pkg/limit"
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/server"
    "github.com/myuser/user-service/kitex_gen/user/userservice"
)

func main() {
    // 服务地址
    addr, _ := net.ResolveTCPAddr("tcp", ":8888")
    
    // 创建服务器
    svr := userservice.NewServer(
        NewUserService(userRepo, cache),
        
        // 服务地址
        server.WithServiceAddr(addr),
        
        // 服务名称
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
            ServiceName: "user-service",
        }),
        
        // 限流配置
        server.WithLimit(&limit.Option{
            MaxConnections: 10000,
            MaxQPS:         5000,
        }),
        
        // 超时配置
        server.WithReadWriteTimeout(5 * time.Second),
        
        // 中间件
        server.WithMiddleware(LogMiddleware),
        server.WithMiddleware(MetricsMiddleware),
        
        // 服务注册
        server.WithRegistry(registry),
        
        // 优雅退出
        server.WithExitWaitTime(10 * time.Second),
    )
    
    klog.Info("Server starting on :8888")
    
    err := svr.Run()
    if err != nil {
        klog.Fatal(err)
    }
}
```

---

## 客户端开发

### 创建客户端

#### 基本用法

```go
package main

import (
    "context"
    "log"
    
    "github.com/cloudwego/kitex/client"
    "github.com/myuser/user-service/kitex_gen/user"
    "github.com/myuser/user-service/kitex_gen/user/userservice"
)

func main() {
    // 创建客户端
    cli, err := userservice.NewClient("user-service", client.WithHostPorts("127.0.0.1:8888"))
    if err != nil {
        log.Fatal(err)
    }
    
    // 调用方法
    ctx := context.Background()
    req := &user.GetUserRequest{UserId: 1}
    
    resp, err := cli.GetUser(ctx, req)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Response: %+v", resp)
}
```

#### 高级配置

```go
package main

import (
    "context"
    "time"
    
    "github.com/cloudwego/kitex/client"
    "github.com/cloudwego/kitex/pkg/retry"
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/myuser/user-service/kitex_gen/user/userservice"
)

func main() {
    // 创建客户端
    cli, err := userservice.NewClient(
        "user-service",
        
        // 服务地址
        client.WithHostPorts("127.0.0.1:8888"),
        
        // 客户端信息
        client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
            ServiceName: "api-gateway",
        }),
        
        // 超时配置
        client.WithRPCTimeout(3 * time.Second),
        client.WithConnectTimeout(1 * time.Second),
        
        // 重试策略
        client.WithRetryPolicy(retry.BuildFailurePolicy(retry.NewFailurePolicy())),
        
        // 连接池配置
        client.WithLongConnection(client.ConnPoolConfig{
            MaxIdlePerAddress: 10,
            MaxIdleGlobal:     100,
            MaxIdleTimeout:    60 * time.Second,
        }),
        
        // 中间件
        client.WithMiddleware(ClientLogMiddleware),
        
        // 服务发现
        client.WithResolver(resolver),
        
        // 负载均衡
        client.WithLoadBalancer(loadbalancer),
    )
    
    if err != nil {
        panic(err)
    }
    
    // 使用客户端
    ctx := context.Background()
    resp, err := cli.GetUser(ctx, &user.GetUserRequest{UserId: 1})
}
```

### 调用选项

```go
package main

import (
    "context"
    "time"
    
    "github.com/cloudwego/kitex/client"
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/myuser/user-service/kitex_gen/user"
)

func CallWithOptions(cli userservice.Client) {
    ctx := context.Background()
    req := &user.GetUserRequest{UserId: 1}
    
    // 为单次调用设置选项
    resp, err := cli.GetUser(
        ctx,
        req,
        // 设置超时
        client.WithRPCTimeout(5*time.Second),
        
        // 设置目标地址
        client.WithHostPort("192.168.1.100:8888"),
        
        // 设置请求标签
        client.WithTag("key", "value"),
    )
}
```

---

## 中间件

### 服务端中间件

```go
package middleware

import (
    "context"
    "time"
    
    "github.com/cloudwego/kitex/pkg/endpoint"
    "github.com/cloudwego/kitex/pkg/klog"
    "github.com/cloudwego/kitex/pkg/rpcinfo"
)

// 日志中间件
func LogMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
    return func(ctx context.Context, req, resp interface{}) error {
        start := time.Now()
        
        // 获取 RPC 信息
        ri := rpcinfo.GetRPCInfo(ctx)
        
        klog.Infof("Request: service=%s, method=%s, from=%s",
            ri.To().ServiceName(),
            ri.To().Method(),
            ri.From().Address(),
        )
        
        // 调用下一个处理器
        err := next(ctx, req, resp)
        
        // 记录响应
        duration := time.Since(start)
        klog.Infof("Response: service=%s, method=%s, duration=%v, err=%v",
            ri.To().ServiceName(),
            ri.To().Method(),
            duration,
            err,
        )
        
        return err
    }
}

// 认证中间件
func AuthMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
    return func(ctx context.Context, req, resp interface{}) error {
        // 从元信息获取 token
        token, ok := metainfo.GetValue(ctx, "auth_token")
        if !ok {
            return errors.New("missing auth token")
        }
        
        // 验证 token
        userID, err := validateToken(token)
        if err != nil {
            return errors.New("invalid token")
        }
        
        // 存储用户 ID
        ctx = metainfo.WithValue(ctx, "user_id", userID)
        
        return next(ctx, req, resp)
    }
}

// 监控中间件
func MetricsMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
    return func(ctx context.Context, req, resp interface{}) error {
        start := time.Now()
        
        ri := rpcinfo.GetRPCInfo(ctx)
        serviceName := ri.To().ServiceName()
        methodName := ri.To().Method()
        
        // 调用下一个处理器
        err := next(ctx, req, resp)
        
        // 记录指标
        duration := time.Since(start)
        
        // 发送到监控系统
        metrics.RecordRPCDuration(serviceName, methodName, duration)
        
        if err != nil {
            metrics.IncrementRPCError(serviceName, methodName)
        } else {
            metrics.IncrementRPCSuccess(serviceName, methodName)
        }
        
        return err
    }
}

// 限流中间件
func RateLimitMiddleware(limiter *rate.Limiter) endpoint.Middleware {
    return func(next endpoint.Endpoint) endpoint.Endpoint {
        return func(ctx context.Context, req, resp interface{}) error {
            if !limiter.Allow() {
                return errors.New("rate limit exceeded")
            }
            return next(ctx, req, resp)
        }
    }
}

func validateToken(token string) (string, error) {
    // 实现 token 验证逻辑
    return "user123", nil
}

type metrics struct{}

func (m *metrics) RecordRPCDuration(service, method string, duration time.Duration) {}
func (m *metrics) IncrementRPCError(service, method string)                          {}
func (m *metrics) IncrementRPCSuccess(service, method string)                        {}
```

### 客户端中间件

```go
package middleware

import (
    "context"
    "time"
    
    "github.com/cloudwego/kitex/pkg/endpoint"
    "github.com/cloudwego/kitex/pkg/klog"
    "github.com/cloudwego/kitex/pkg/metainfo"
    "github.com/cloudwego/kitex/pkg/rpcinfo"
)

// 客户端日志中间件
func ClientLogMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
    return func(ctx context.Context, req, resp interface{}) error {
        start := time.Now()
        
        ri := rpcinfo.GetRPCInfo(ctx)
        
        klog.Infof("Client request: service=%s, method=%s",
            ri.To().ServiceName(),
            ri.To().Method(),
        )
        
        err := next(ctx, req, resp)
        
        duration := time.Since(start)
        klog.Infof("Client response: service=%s, method=%s, duration=%v, err=%v",
            ri.To().ServiceName(),
            ri.To().Method(),
            duration,
            err,
        )
        
        return err
    }
}

// 客户端元信息中间件
func ClientMetainfoMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
    return func(ctx context.Context, req, resp interface{}) error {
        // 添加元信息
        ctx = metainfo.WithPersistentValue(ctx, "request_id", generateRequestID())
        ctx = metainfo.WithPersistentValue(ctx, "timestamp", time.Now().Format(time.RFC3339))
        
        return next(ctx, req, resp)
    }
}

// 客户端重试中间件
func ClientRetryMiddleware(maxRetries int) endpoint.Middleware {
    return func(next endpoint.Endpoint) endpoint.Endpoint {
        return func(ctx context.Context, req, resp interface{}) error {
            var err error
            for i := 0; i <= maxRetries; i++ {
                err = next(ctx, req, resp)
                if err == nil {
                    return nil
                }
                
                klog.Warnf("Request failed (attempt %d/%d): %v", i+1, maxRetries+1, err)
                
                if i < maxRetries {
                    time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
                }
            }
            return err
        }
    }
}

func generateRequestID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}
```

---

## 服务治理

### 服务注册与发现

#### 使用 Etcd

```go
package main

import (
    "log"
    
    "github.com/cloudwego/kitex/client"
    "github.com/cloudwego/kitex/pkg/registry"
    "github.com/cloudwego/kitex/server"
    etcd "github.com/kitex-contrib/registry-etcd"
    "github.com/myuser/user-service/kitex_gen/user/userservice"
)

// 服务端注册
func StartServer() {
    // 创建 Etcd 注册中心
    r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
    if err != nil {
        log.Fatal(err)
    }
    
    // 创建服务器
    svr := userservice.NewServer(
        NewUserService(),
        server.WithServiceAddr(&net.TCPAddr{Port: 8888}),
        server.WithRegistry(r),
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
            ServiceName: "user-service",
        }),
    )
    
    err = svr.Run()
    if err != nil {
        log.Fatal(err)
    }
}

// 客户端发现
func CreateClient() {
    // 创建 Etcd 解析器
    r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
    if err != nil {
        log.Fatal(err)
    }
    
    // 创建客户端
    cli, err := userservice.NewClient(
        "user-service",
        client.WithResolver(r),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // 使用客户端
    ctx := context.Background()
    resp, err := cli.GetUser(ctx, &user.GetUserRequest{UserId: 1})
}
```

#### 使用 Consul

```go
import (
    consul "github.com/kitex-contrib/registry-consul"
)

// 服务端
func StartServer() {
    r, err := consul.NewConsulRegister("127.0.0.1:8500")
    if err != nil {
        log.Fatal(err)
    }
    
    svr := userservice.NewServer(
        NewUserService(),
        server.WithRegistry(r),
    )
    
    svr.Run()
}

// 客户端
func CreateClient() {
    r, err := consul.NewConsulResolver("127.0.0.1:8500")
    if err != nil {
        log.Fatal(err)
    }
    
    cli, err := userservice.NewClient(
        "user-service",
        client.WithResolver(r),
    )
}
```

### 负载均衡

```go
package main

import (
    "github.com/cloudwego/kitex/client"
    "github.com/cloudwego/kitex/pkg/loadbalance"
)

func CreateClientWithLB() {
    cli, err := userservice.NewClient(
        "user-service",
        // 轮询
        client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()),
        
        // 或随机
        // client.WithLoadBalancer(loadbalance.NewWeightedRandomBalancer()),
        
        // 或一致性哈希
        // client.WithLoadBalancer(loadbalance.NewConsistHashBalancer()),
    )
}
```

### 熔断

```go
package main

import (
    "github.com/cloudwego/kitex/client"
    "github.com/cloudwego/kitex/pkg/circuitbreak"
)

func CreateClientWithCircuitBreaker() {
    // 熔断器配置
    cbSuite := circuitbreak.NewCBSuite(func(ri rpcinfo.RPCInfo) string {
        return circuitbreak.RPCInfo2Key(ri)
    })
    
    cli, err := userservice.NewClient(
        "user-service",
        client.WithCircuitBreaker(cbSuite),
    )
}
```

### 限流

```go
package main

import (
    "github.com/cloudwego/kitex/pkg/limit"
    "github.com/cloudwego/kitex/server"
)

func StartServerWithLimit() {
    svr := userservice.NewServer(
        NewUserService(),
        // QPS 限流
        server.WithLimit(&limit.Option{
            MaxConnections: 10000,  // 最大连接数
            MaxQPS:         5000,   // 最大 QPS
        }),
    )
    
    svr.Run()
}
```

---

## 性能优化

### 1. 使用连接池

```go
func CreateClient() {
    cli, err := userservice.NewClient(
        "user-service",
        client.WithLongConnection(client.ConnPoolConfig{
            MaxIdlePerAddress: 10,              // 每个地址最大空闲连接数
            MaxIdleGlobal:     100,             // 全局最大空闲连接数
            MaxIdleTimeout:    60 * time.Second, // 空闲连接超时
        }),
    )
}
```

### 2. 使用 Netpoll

```go
import (
    "github.com/cloudwego/kitex/pkg/remote/trans/netpoll"
    "github.com/cloudwego/kitex/server"
)

func StartServer() {
    svr := userservice.NewServer(
        NewUserService(),
        // 使用 netpoll
        server.WithTransHandlerFactory(netpoll.NewSvrTransHandlerFactory()),
    )
    
    svr.Run()
}
```

### 3. 对象池

```go
var requestPool = sync.Pool{
    New: func() interface{} {
        return &user.GetUserRequest{}
    },
}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
    // 使用对象池
    pooledReq := requestPool.Get().(*user.GetUserRequest)
    defer requestPool.Put(pooledReq)
    
    // 复制请求数据
    *pooledReq = *req
    
    // 处理逻辑
    return &user.GetUserResponse{}, nil
}
```

### 4. 批量请求

```go
// 批量获取用户
func (s *UserServiceImpl) BatchGetUsers(ctx context.Context, req *user.BatchGetUsersRequest) (*user.BatchGetUsersResponse, error) {
    // 批量查询数据库
    users, err := s.userRepo.BatchGet(req.UserIds)
    if err != nil {
        return nil, err
    }
    
    // 转换结果
    thriftUsers := make([]*user.User, 0, len(users))
    for _, u := range users {
        thriftUsers = append(thriftUsers, convertToThriftUser(u))
    }
    
    return &user.BatchGetUsersResponse{
        Code:    0,
        Message: "success",
        Users:   thriftUsers,
    }, nil
}
```

---

## 实战案例

### 完整的微服务示例

#### 项目结构

```
user-service/
├── idl/
│   └── user.thrift
├── kitex_gen/
│   └── user/
├── handler/
│   └── user_handler.go
├── repository/
│   └── user_repository.go
├── middleware/
│   ├── log.go
│   └── auth.go
├── config/
│   └── config.go
├── main.go
├── go.mod
└── go.sum
```

#### IDL 定义

**idl/user.thrift**
```thrift
namespace go user

struct User {
    1: required i64 id
    2: required string username
    3: required string email
    4: optional i32 age
    5: optional string created_at
    6: optional string updated_at
}

struct CreateUserRequest {
    1: required string username
    2: required string email
    3: optional i32 age
}

struct CreateUserResponse {
    1: required i32 code
    2: required string message
    3: optional User user
}

struct GetUserRequest {
    1: required i64 user_id
}

struct GetUserResponse {
    1: required i32 code
    2: required string message
    3: optional User user
}

struct UpdateUserRequest {
    1: required i64 user_id
    2: optional string username
    3: optional string email
    4: optional i32 age
}

struct UpdateUserResponse {
    1: required i32 code
    2: required string message
    3: optional User user
}

struct DeleteUserRequest {
    1: required i64 user_id
}

struct DeleteUserResponse {
    1: required i32 code
    2: required string message
}

struct ListUsersRequest {
    1: optional i32 page = 1
    2: optional i32 page_size = 10
}

struct ListUsersResponse {
    1: required i32 code
    2: required string message
    3: optional list<User> users
    4: optional i64 total
}

service UserService {
    CreateUserResponse CreateUser(1: CreateUserRequest req)
    GetUserResponse GetUser(1: GetUserRequest req)
    UpdateUserResponse UpdateUser(1: UpdateUserRequest req)
    DeleteUserResponse DeleteUser(1: DeleteUserRequest req)
    ListUsersResponse ListUsers(1: ListUsersRequest req)
}
```

#### 生成代码

```bash
kitex -module github.com/myuser/user-service \
    -service userservice \
    idl/user.thrift

go mod tidy
```

#### Handler 实现

**handler/user_handler.go**
```go
package handler

import (
    "context"
    "time"
    
    "github.com/myuser/user-service/kitex_gen/user"
    "github.com/myuser/user-service/repository"
)

type UserServiceImpl struct {
    repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserServiceImpl {
    return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
    if req.Username == "" || req.Email == "" {
        return &user.CreateUserResponse{
            Code:    400,
            Message: "username and email are required",
        }, nil
    }
    
    newUser := &repository.User{
        Username:  req.Username,
        Email:     req.Email,
        Age:       int(req.GetAge()),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := s.repo.Create(newUser); err != nil {
        return &user.CreateUserResponse{
            Code:    500,
            Message: "failed to create user",
        }, err
    }
    
    return &user.CreateUserResponse{
        Code:    0,
        Message: "success",
        User:    convertToThriftUser(newUser),
    }, nil
}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
    u, err := s.repo.GetByID(req.UserId)
    if err != nil {
        if err == repository.ErrUserNotFound {
            return &user.GetUserResponse{
                Code:    404,
                Message: "user not found",
            }, nil
        }
        return &user.GetUserResponse{
            Code:    500,
            Message: "internal error",
        }, err
    }
    
    return &user.GetUserResponse{
        Code:    0,
        Message: "success",
        User:    convertToThriftUser(u),
    }, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
    u, err := s.repo.GetByID(req.UserId)
    if err != nil {
        if err == repository.ErrUserNotFound {
            return &user.UpdateUserResponse{
                Code:    404,
                Message: "user not found",
            }, nil
        }
        return &user.UpdateUserResponse{
            Code:    500,
            Message: "internal error",
        }, err
    }
    
    if req.IsSetUsername() {
        u.Username = req.GetUsername()
    }
    if req.IsSetEmail() {
        u.Email = req.GetEmail()
    }
    if req.IsSetAge() {
        u.Age = int(req.GetAge())
    }
    u.UpdatedAt = time.Now()
    
    if err := s.repo.Update(u); err != nil {
        return &user.UpdateUserResponse{
            Code:    500,
            Message: "failed to update user",
        }, err
    }
    
    return &user.UpdateUserResponse{
        Code:    0,
        Message: "success",
        User:    convertToThriftUser(u),
    }, nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
    if err := s.repo.Delete(req.UserId); err != nil {
        if err == repository.ErrUserNotFound {
            return &user.DeleteUserResponse{
                Code:    404,
                Message: "user not found",
            }, nil
        }
        return &user.DeleteUserResponse{
            Code:    500,
            Message: "failed to delete user",
        }, err
    }
    
    return &user.DeleteUserResponse{
        Code:    0,
        Message: "success",
    }, nil
}

func (s *UserServiceImpl) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
    page := int(req.GetPage())
    pageSize := int(req.GetPageSize())
    
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 10
    }
    
    users, total, err := s.repo.List(page, pageSize)
    if err != nil {
        return &user.ListUsersResponse{
            Code:    500,
            Message: "failed to list users",
        }, err
    }
    
    thriftUsers := make([]*user.User, 0, len(users))
    for _, u := range users {
        thriftUsers = append(thriftUsers, convertToThriftUser(u))
    }
    
    return &user.ListUsersResponse{
        Code:    0,
        Message: "success",
        Users:   thriftUsers,
        Total:   &total,
    }, nil
}

func convertToThriftUser(u *repository.User) *user.User {
    return &user.User{
        Id:        u.ID,
        Username:  u.Username,
        Email:     u.Email,
        Age:       func() *int32 { a := int32(u.Age); return &a }(),
        CreatedAt: func() *string { s := u.CreatedAt.Format(time.RFC3339); return &s }(),
        UpdatedAt: func() *string { s := u.UpdatedAt.Format(time.RFC3339); return &s }(),
    }
}
```

#### Repository

**repository/user_repository.go**
```go
package repository

import (
    "errors"
    "sync"
    "time"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
    ID        int64
    Username  string
    Email     string
    Age       int
    CreatedAt time.Time
    UpdatedAt time.Time
}

type UserRepository struct {
    mu     sync.RWMutex
    users  map[int64]*User
    nextID int64
}

func NewUserRepository() *UserRepository {
    return &UserRepository{
        users:  make(map[int64]*User),
        nextID: 1,
    }
}

func (r *UserRepository) Create(u *User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    u.ID = r.nextID
    r.users[u.ID] = u
    r.nextID++
    
    return nil
}

func (r *UserRepository) GetByID(id int64) (*User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    u, ok := r.users[id]
    if !ok {
        return nil, ErrUserNotFound
    }
    
    return u, nil
}

func (r *UserRepository) Update(u *User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, ok := r.users[u.ID]; !ok {
        return ErrUserNotFound
    }
    
    r.users[u.ID] = u
    return nil
}

func (r *UserRepository) Delete(id int64) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, ok := r.users[id]; !ok {
        return ErrUserNotFound
    }
    
    delete(r.users, id)
    return nil
}

func (r *UserRepository) List(page, pageSize int) ([]*User, int64, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    total := int64(len(r.users))
    
    offset := (page - 1) * pageSize
    limit := pageSize
    
    users := make([]*User, 0, limit)
    count := 0
    for _, u := range r.users {
        if count >= offset && len(users) < limit {
            users = append(users, u)
        }
        count++
    }
    
    return users, total, nil
}
```

#### 主函数

**main.go**
```go
package main

import (
    "log"
    "net"
    
    "github.com/cloudwego/kitex/pkg/klog"
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/server"
    "github.com/myuser/user-service/handler"
    "github.com/myuser/user-service/kitex_gen/user/userservice"
    "github.com/myuser/user-service/middleware"
    "github.com/myuser/user-service/repository"
)

func main() {
    // 初始化 repository
    repo := repository.NewUserRepository()
    
    // 创建 handler
    userHandler := handler.NewUserService(repo)
    
    // 服务地址
    addr, _ := net.ResolveTCPAddr("tcp", ":8888")
    
    // 创建服务器
    svr := userservice.NewServer(
        userHandler,
        server.WithServiceAddr(addr),
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
            ServiceName: "user-service",
        }),
        server.WithMiddleware(middleware.LogMiddleware),
        server.WithMiddleware(middleware.MetricsMiddleware),
    )
    
    klog.Info("User service starting on :8888")
    
    err := svr.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

#### 客户端示例

**client/main.go**
```go
package main

import (
    "context"
    "log"
    
    "github.com/cloudwego/kitex/client"
    "github.com/myuser/user-service/kitex_gen/user"
    "github.com/myuser/user-service/kitex_gen/user/userservice"
)

func main() {
    // 创建客户端
    cli, err := userservice.NewClient("user-service", client.WithHostPorts("127.0.0.1:8888"))
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 创建用户
    createResp, err := cli.CreateUser(ctx, &user.CreateUserRequest{
        Username: "zhangsan",
        Email:    "zhangsan@example.com",
        Age:      func() *int32 { a := int32(25); return &a }(),
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Create user response: %+v", createResp)
    
    // 获取用户
    userID := createResp.User.Id
    getResp, err := cli.GetUser(ctx, &user.GetUserRequest{
        UserId: userID,
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Get user response: %+v", getResp)
    
    // 更新用户
    updateResp, err := cli.UpdateUser(ctx, &user.UpdateUserRequest{
        UserId:   userID,
        Username: func() *string { s := "zhangsan_updated"; return &s }(),
        Age:      func() *int32 { a := int32(26); return &a }(),
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Update user response: %+v", updateResp)
    
    // 列出用户
    listResp, err := cli.ListUsers(ctx, &user.ListUsersRequest{
        Page:     func() *int32 { p := int32(1); return &p }(),
        PageSize: func() *int32 { s := int32(10); return &s }(),
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("List users response: %+v", listResp)
    
    // 删除用户
    deleteResp, err := cli.DeleteUser(ctx, &user.DeleteUserRequest{
        UserId: userID,
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Delete user response: %+v", deleteResp)
}
```

---

## 总结

本文详细介绍了 Kitex 框架的各个方面：

1. **CLI 工具**：kitex 命令行工具的完整使用方法和参数详解
2. **Context**：context.Context 和 metainfo 的使用方法
3. **IDL 定义**：Thrift 和 Protobuf IDL 编写规范
4. **服务端**：服务实现、启动配置
5. **客户端**：客户端创建、调用方法
6. **中间件**：服务端和客户端中间件开发
7. **服务治理**：注册发现、负载均衡、熔断、限流
8. **性能优化**：各种性能优化技巧
9. **实战**：完整的微服务示例

### 最佳实践

1. 使用 kitex 工具生成代码
2. 合理设计 IDL 接口
3. 使用中间件处理通用逻辑
4. 使用服务注册发现实现动态路由
5. 配置合理的超时和重试策略
6. 使用连接池提升性能
7. 集成监控和链路追踪
8. 完善的错误处理
9. 使用元信息传递上下文
10. 编写单元测试和集成测试

更多详细信息请参考 [Kitex 官方文档](https://www.cloudwego.io/zh/docs/kitex/)。
